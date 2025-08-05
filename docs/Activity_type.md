Структура fit.ActivityFile

```go
// ActivityFile - основная структура для представления данных спортивной активности из FIT-файла
type ActivityFile struct {
	Activity    *ActivityMsg
	Sessions    []*SessionMsg
	Laps        []*LapMsg
	Lengths     []*LengthMsg
	Records     []*RecordMsg
	Events      []*EventMsg
	Hrvs        []*HrvMsg
	DeviceInfos []*DeviceInfoMsg
}
```

