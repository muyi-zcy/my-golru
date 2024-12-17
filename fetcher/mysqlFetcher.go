package fetcher

type DataFetcher interface {
	FetchData(key string) (string, error)
}
