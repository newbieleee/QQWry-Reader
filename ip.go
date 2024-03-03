package ip

import (
    `bufio`
    `encoding/binary`
    `errors`
    `io`
    `net`

    `golang.org/x/text/encoding/simplifiedchinese`
)

const (
    indexOffsetSize  int64 = 4
    ipSize           int64 = 4
    offsetSize       int64 = 3
    indexSize        int64 = 7
    redirectFlagSize int64 = 1
    redirectMixed    byte  = 1 // 该模式指向offset或字符串
    redirectString   byte  = 2 // 该模式一定指向字符串
)

type Source interface {
    Open() error
    io.Reader
    io.ReaderAt
    io.Seeker
    io.Closer
}

type Data struct {
    source  Source
    meta    indexMeta
    version string
}

type indexMeta struct {
    begin uint32
    end   uint32
    total uint32
}

type indexRecord struct {
    ipBegin uint32
    offset  int64
}

type Record struct {
    Country  []byte
    Region   []byte
    Location []byte
}

func (i indexMeta) calcOffset(index int64) int64 {
    return int64(i.begin) + index*indexSize
}

var (
    VersionIPBegin uint32 = 4294967040
    VersionIPEnd   uint32 = 4294967295

    UnknownText          = []byte("未知")
    DefaultUnknownRegion = []byte{32, 67, 90, 56, 56, 46, 78, 69, 84} // 未知地区为[ CZ88.NET]

    ErrIP       = errors.New("IP不正确")
    ErrNotFound = errors.New("未找到IP记录")
    ErrDecode   = errors.New("已找到记录，编码解析错误")
)

func (data *Data) Query(ip string) (Record, error) {
    ipBytes := net.ParseIP(ip)
    if ipBytes == nil {
        return Record{}, ErrIP
    }
    ipInt := binary.BigEndian.Uint32(ipBytes.To4())
    if ipInt >= VersionIPBegin && ipInt <= VersionIPEnd {
        r := Record{Country: []byte("保留地址")}
        r.Location = r.Country
        return r, nil
    }
    var (
        left  int64
        right = int64(data.meta.total - 1)
    )
    for left <= right {
        mid := left + ((right - left) >> 1)
        off := data.meta.calcOffset(mid)
        midRecord, err := data.readIndex(off)
        if err != nil {
            return Record{}, err
        }
        switch {
        case midRecord.ipBegin > ipInt:
            right = mid - 1
        case midRecord.ipBegin == ipInt:
            return data.readRecord(midRecord.offset + ipSize)
        default:
            ipEnd, err := data.readIP(midRecord.offset)
            if err != nil {
                return Record{}, err
            }
            if ipEnd >= ipInt {
                return data.readRecord(midRecord.offset + ipSize)
            }
            left = mid + 1
        }
    }
    return Record{}, ErrNotFound
}

func (data *Data) readIP(pos int64) (uint32, error) {
    buf := make([]byte, ipSize)
    _, err := data.source.ReadAt(buf, pos)
    if err != nil {
        return 0, err
    }
    return binary.LittleEndian.Uint32(buf), nil
}

func (data *Data) readRecord(pos int64) (res Record, err error) {
    countryGBK, off, err := data.readString(pos)
    if err != nil {
        return
    }
    regionGBK, off, err := data.readString(off)
    if err != nil {
        return
    }
    decoder := simplifiedchinese.GBK.NewDecoder()
    country, err := decoder.Bytes(countryGBK[:len(countryGBK)-1])
    if err != nil {
        return
    }
    region, err := decoder.Bytes(regionGBK[:len(regionGBK)-1])
    if err != nil {
        err = ErrDecode
        return
    }
    region = replaceUnknownChars(region)
    location := make([]byte, len(country)+len(region))
    n := copy(location, country)
    copy(location[n:], region)
    return Record{country, region, location}, nil
}

func replaceUnknownChars(chars []byte) []byte {
    if len(chars) != len(DefaultUnknownRegion) {
        return chars
    }
    n := 0
    for i := range chars {
        if chars[i] == DefaultUnknownRegion[i] {
            n++
        }
    }
    if n == len(DefaultUnknownRegion) {
        return UnknownText
    }
    return chars
}

func (data *Data) readString(pos int64) (result []byte, nextOffset int64, err error) {
    for {
        var (
            m   byte
            off int64
        )
        m, off, err = data.readMode(pos)
        if err != nil {
            return
        }
        if off == 0 && (m == redirectMixed || m == redirectString) {
            result = UnknownText
            if nextOffset == 0 {
                nextOffset = pos + redirectFlagSize + offsetSize
            }
            return
        }
        switch m {
        case redirectMixed:
            pos = off
        case redirectString:
            nextOffset = pos + redirectFlagSize + offsetSize
            pos = off
        default:
            // _, err = data.source.Seek(pos, io.SeekStart)
            // if err != nil {
            //     return
            // }
            result, err = bufio.NewReader(&concurrencyReader{reader: data.source, pos: pos}).ReadBytes(0)
            if err != nil {
                return
            }
            if nextOffset == 0 {
                nextOffset = pos + int64(len(result))
            }
            return
        }
    }
}

func (data *Data) readMode(pos int64) (mode byte, offset int64, err error) {
    buf := make([]byte, redirectFlagSize+offsetSize)
    _, err = data.source.ReadAt(buf, pos)
    if err != nil {
        return
    }
    return buf[0], int64(binary.LittleEndian.Uint32(buf) >> 8), nil
}

func (data *Data) readIndex(pos int64) (indexRecord, error) {
    buf := make([]byte, ipSize+offsetSize)
    _, err := data.source.ReadAt(buf, pos)
    if err != nil {
        return indexRecord{}, err
    }
    return indexRecord{
        ipBegin: binary.LittleEndian.Uint32(buf[:ipSize]),
        offset:  int64(binary.LittleEndian.Uint32(buf[offsetSize:]) >> 8),
    }, nil
}

func (data *Data) loadSource() error {
    err := data.source.Open()
    if err != nil {
        return err
    }
    buf := make([]byte, indexOffsetSize*2)
    _, err = data.source.Read(buf)
    if err != nil {
        return err
    }
    data.meta = indexMeta{
        begin: binary.LittleEndian.Uint32(buf[:indexOffsetSize]),
        end:   binary.LittleEndian.Uint32(buf[indexOffsetSize:]),
    }
    data.meta.total = (data.meta.end - data.meta.begin) / uint32(indexSize)
    return nil
}

func (data *Data) loadVersion() error {
    index, err := data.readIndex(int64(data.meta.end))
    if err != nil {
        return err
    }
    r, err := data.readRecord(index.offset)
    if err != nil {
        return err
    }
    v := make([]byte, 8)
    off := 0
    for _, b := range r.Region {
        if b >= 48 && b <= 57 {
            v[off] = b
            off++
        }
    }
    data.version = string(v[:off])
    return nil
}

func (data *Data) Version() string {
    return data.version
}

func (data *Data) Close() error {
    return data.source.Close()
}

func New(builder Builder) (*Data, error) {
    data := &Data{}
    err := builder.build(data)
    if err != nil {
        return nil, err
    }
    err = data.loadSource()
    if err != nil {
        return nil, err
    }
    err = data.loadVersion()
    if err != nil {
        return nil, err
    }
    return data, nil
}
