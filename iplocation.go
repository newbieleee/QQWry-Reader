package ip

import (
    `bufio`
    `encoding/binary`
    `errors`
    `net`
    `os`
    `path/filepath`

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

type Data struct {
    FilePath string
    source   *os.File
    meta     indexMeta
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

type IPRecord struct {
    country []byte
    region  []byte
}

func (i indexMeta) calcOffset(index int64) int64 {
    return int64(i.begin) + index*indexSize
}

var (
    UnknownText          = []byte("未知")
    DefaultUnknownRegion = []byte{67, 90, 56, 56, 46, 78, 69, 84} // 未知地区为CZ88.NET

    ErrIP       = errors.New("IP不正确")
    ErrNotFound = errors.New("未找到IP记录")
    ErrDecode   = errors.New("已找到记录，转码出现错误")
)

func (data *Data) Query(ip string) (IPRecord, error) {
    ipBytes := net.ParseIP(ip)
    if ipBytes == nil {
        return IPRecord{}, ErrIP
    }
    err := data.loadFile()
    if err != nil {
        return IPRecord{}, err
    }
    defer func() {
        _ = data.source.Close()
    }()
    ipInt := binary.BigEndian.Uint32(ipBytes.To4())
    var (
        left  int64
        right = int64(data.meta.total - 1)
    )
    for left <= right {
        mid := left + ((right - left) >> 1)
        off := data.meta.calcOffset(mid)
        midRecord, err := data.readIndex(off)
        if err != nil {
            return IPRecord{}, err
        }
        switch {
        case midRecord.ipBegin > ipInt:
            right = mid + 1
        case midRecord.ipBegin == ipInt:
            return data.readRecord(midRecord.offset + ipSize)
        default:
            ipEnd, err := data.readIP(midRecord.offset)
            if err != nil {
                return IPRecord{}, err
            }
            if ipEnd >= ipInt {
                return data.readRecord(midRecord.offset + ipSize)
            }
            left = mid + 1
        }
    }
    return IPRecord{}, ErrNotFound
}

func (data *Data) readIP(pos int64) (uint32, error) {
    buf := make([]byte, ipSize)
    _, err := data.source.ReadAt(buf, pos)
    if err != nil {
        return 0, err
    }
    return binary.LittleEndian.Uint32(buf), nil
}

func (data *Data) readRecord(pos int64) (res IPRecord, err error) {
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
    return IPRecord{country, replaceUnknownChars(region)}, nil
}

func replaceUnknownChars(chars []byte) []byte {
    if len(chars) != 8 {
        return chars
    }
    for i := range chars {
        if chars[i] != DefaultUnknownRegion[i] {
            return UnknownText
        }
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
            _, err = data.source.Seek(pos, 0)
            if err != nil {
                return
            }
            result, err = bufio.NewReader(data.source).ReadBytes(0)
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

func (data *Data) loadFile() error {
    f, err := os.OpenFile(data.FilePath, os.O_RDONLY, 0644)
    if err != nil {
        return err
    }
    buf := make([]byte, indexOffsetSize*2)
    _, err = f.Read(buf)
    if err != nil {
        return err
    }
    data.source = f
    data.meta = indexMeta{
        begin: binary.LittleEndian.Uint32(buf[:indexOffsetSize]),
        end:   binary.LittleEndian.Uint32(buf[indexOffsetSize:]),
    }
    data.meta.total = (data.meta.end - data.meta.begin) / uint32(indexSize)
    return nil
}

func (data *Data) ReplaceUnknownFlag() {

}

func New(filename string) (*Data, error) {
    path := filepath.Clean(filename)
    path, err := filepath.Abs(path)
    if err != nil {
        return nil, err
    }
    return &Data{FilePath: path}, nil
}
