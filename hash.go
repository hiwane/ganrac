package ganrac

type Hash uint64

type Hashable interface {
	Hash() Hash
	equaler
}
