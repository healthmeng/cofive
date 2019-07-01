package ai

type MTTree struct{
	x,y []int
	win,all int
	nexts []MTTree
	hash []byte
	pathname string
}

func SearchStep(player* AIPlayer) *MTTree{
	return nil
}

func (mt *MTTree)LoadAllNext() int{
	// return nSteps
	// if not exists, create them in db
	return 0
}

func (mt *MTTree)GetMaxRatioStep() *MTTree{
	if len(mt.nexts)==0{
		if mt.LoadAllNext()<=0{
			return nil
		}
	}

	ratio:=0.0
	n:=&mt.nexts[0]
	for _,s:=range(mt.nexts){
		r:=float64(s.win)/float64(s.all)
		if r>ratio{
			ratio=r
			n=&s
		}
	}
	return n
}
