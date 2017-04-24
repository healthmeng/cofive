#include<stdio.h>
#include <assert.h>
#include <time.h>
#include <stdlib.h>
#include "robot.h"

#define SCORE_INIT -20000


static int frame[15][15];
static curstep=0;
static int max_score;
static ailevel=0;
static struct STEP steplog[MAX_STEP];

typedef struct _BWSCORE{
	int bscore;
	int wscore;
}BWCORE;

int computerbw=2;	// computer use white

void addcond(BWCORE *ps, int bw,int cnt,int left,int right, int midspace, int nturn){
	int score=0;
	if(midspace==cnt)
		midspace=0;
	if( cnt>=5 ){
		if (midspace>=5 || cnt-midspace>=5)
			score=50000;
		else score=720;
	}
	if( cnt==4){
		if(!midspace){
			if( left && right)
			{
				score=4320;
			}
			else if(left || right)
			{
				score=720;
			}
		}else{
				score=720;
		}
		if(score && bw==nturn)
			score=10000;
	}
	if (cnt==3){
		if(left && right && left+right>2)
		{
			score=720;
			if (bw==nturn) score=2000;
		}
		else{
		 if(left+right>=2)
			score=120;
		 if(midspace)
			score-=40;
		}
	}

	if (cnt==2){
		if(left && right && left+right>3)
			score=120;
		else if(left +right>=3)
			score=20;
		if( score && bw==nturn) score+=10;
	}

	if(cnt==1){
		if(left && right && right+left>=4)
			score=20;
	}
	if(bw==1)
		ps->bscore+=score;
	else if(bw==2)
		ps->wscore+=score;
//    printf("addcount: bw %d, cnt %d\n",bw,cnt, score);
}

BWCORE evaluate(int nextturn){
	BWCORE s={
		.bscore=0,
		.wscore=0
	};
	// |
	int i,j;
	for(i=0;i<15;++i){
		int incal=0,last=0; // 0,1,2
		int lsp=0,rsp=0,msp=0;
		int cnt=0;
		for(j=0;j<15;++j){
			int cur=frame[i][j];
			if(!incal){
				if(cur==0) // ---
					lsp++;
				else{
					incal=cur;// --*
					cnt=1;
					if(j==14)
						addcond(&s,incal,cnt,lsp,0,0,nextturn);
				}
			}else{ // incal!=0
				if(cur==0){ // -
					if(last==0){ // --
						int k;
						assert(msp!=0); //  --
						if(msp==cnt)
							msp=0;
						for(k=j;k<15 && !frame[i][j];++k)
							rsp++;
						addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						cnt=rsp=incal=0;
						lsp=2;
					}else{ // last!=0 ***-
						if(j==14){
							addcond(&s,incal,cnt,lsp,1,msp,nextturn);
							continue;
						}
						if(msp==0){ 
							msp=cnt;
							rsp=1;
						}else{ // msp!=0 && last!=0: *-**-
							if(frame[i][j+1]==0){ // *-**--
								addcond(&s,incal,cnt,lsp,2,msp,nextturn);
								cnt=incal=rsp=0;
								lsp=1;
							}else{ // *-**-?
								addcond(&s,incal,cnt,lsp,1,msp,nextturn);
								if(frame[i][j+1]==incal){ // *-**-*
									cnt=cnt-msp;
									lsp=1;
									msp=cnt;
									rsp=1;
								}else{ // *--**-x
									cnt=0;
									lsp=1;
									msp=0;
									rsp=0;
									incal=0;
								}
							}
						}
					}

				}else{	// incal!=0, cur!=0 : *x || * x ||  ** ||  *-*
					if(j==14){
						if(cur!=incal)
							addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						else
							addcond(&s,incal,cnt+1,lsp,0,msp,nextturn);
						continue;
					}
					if(incal==cur){
						cnt++;
						rsp=0;
					}else{
						addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						incal=cur;
						lsp=rsp;
						rsp=0;
						msp=0;
						cnt=1;
					}
					
				}
			}
			last=cur;
		}
	}

// -
	for(i=0;i<15;++i){
		int incal=0,last=0; // 0,1,2
		int lsp=0,rsp=0,msp=0;
		int cnt=0;
		for(j=0;j<15;++j){
			int cur=frame[j][i];
			if(!incal){
				if(cur==0) // ---
					lsp++;
				else{
					incal=cur;// --*
					cnt=1;
					if(j==14)
						addcond(&s,incal,cnt,lsp,0,0,nextturn);
				}
			}else{ // incal!=0
				if(cur==0){ // -
					if(last==0){ // --
						int k;
						assert(msp!=0); //  --
						if(msp==cnt)
							msp=0;
						for(k=j;k<15 && !frame[i][j];++k)
							rsp++;
						addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						cnt=rsp=incal=0;
						lsp=2;
					}else{ // last!=0 ***-
						if(j==14){
							addcond(&s,incal,cnt,lsp,1,msp,nextturn);
							continue;
						}
						if(msp==0){ 
							msp=cnt;
							rsp=1;
						}else{ // msp!=0 && last!=0: *-**-
							if(frame[j+1][i]==0){ // *-**--
								addcond(&s,incal,cnt,lsp,2,msp,nextturn);
								cnt=incal=rsp=0;
								lsp=1;
							}else{ // *-**-?
								addcond(&s,incal,cnt,lsp,1,msp,nextturn);
								if(frame[j+1][i]==incal){ // *-**-*
									cnt=cnt-msp;
									lsp=1;
									msp=cnt;
									rsp=1;
								}else{ // *--**-x
									cnt=0;
									lsp=1;
									msp=0;
									rsp=0;
									incal=0;
								}
							}
						}
					}
				}else{	// incal!=0, cur!=0 : *x || * x ||  ** ||  *-*
					if(j==14){
						if(cur!=incal)
							addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						else
							addcond(&s,incal,cnt+1,lsp,0,msp,nextturn);
						continue;
					}
					if(incal==cur){
						cnt++;
						rsp=0;
					}else{
						addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						incal=cur;
						lsp=rsp;
						rsp=0;
						msp=0;
						cnt=1;
					}
					
				}
			}
			last=cur;
		}
	}


// half /
	for(i=4;i<15;++i){
		int incal=0,last=0; // 0,1,2
		int lsp=0,rsp=0,msp=0;
		int cnt=0;
		for(j=0;j<=i;++j){
			int cur=frame[j][i-j];
			if(!incal){
				if(cur==0) // ---
					lsp++;
				else{
					incal=cur;// --*
					cnt=1;
					if(j==i)
						addcond(&s,incal,cnt,lsp,0,0,nextturn);
				}
			}else{ // incal!=0
				if(cur==0){ // -
					if(last==0){ // --
						int k;
						assert(msp!=0); //  --
						if(msp==cnt)
							msp=0;
						for(k=j;k<=i && !frame[k][i-k];++k)
							rsp++;
						addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						cnt=rsp=incal=0;
						lsp=2;
					}else{ // last!=0 ***-
						if(j==i){
							addcond(&s,incal,cnt,lsp,1,msp,nextturn);
							continue;
						}
						if(msp==0){ 
							msp=cnt;
							rsp=1;
						}else{ // msp!=0 && last!=0: *-**-
							if(frame[j+1][i-j-1]==0){ // *-**--
								addcond(&s,incal,cnt,lsp,2,msp,nextturn);
								cnt=incal=rsp=0;
								lsp=1;
							}else{ // *-**-?
								addcond(&s,incal,cnt,lsp,1,msp,nextturn);
								if(frame[j+1][i-j-1]==incal){ // *-**-*
									cnt=cnt-msp;
									lsp=1;
									msp=cnt;
									rsp=1;
								}else{ // *--**-x
									cnt=0;
									lsp=1;
									msp=0;
									rsp=0;
									incal=0;
								}
							}
						}
					}
				}else{	// incal!=0, cur!=0 : *x || * x ||  ** ||  *-*
					if(j==i){
						if(cur!=incal)
							addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						else
							addcond(&s,incal,cnt+1,lsp,0,msp,nextturn);
						continue;
					}
					if(incal==cur){
						cnt++;
						rsp=0;
					}else{
						addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						incal=cur;
						lsp=rsp;
						rsp=0;
						msp=0;
						cnt=1;
					}
					
				}
			}
			last=cur;
		}
	}
// left part of /
	for(i=13;i>=4;--i){
		int incal=0,last=0; // 0,1,2
		int lsp=0,rsp=0,msp=0;
		int cnt=0;
		for(j=0;j<=i;++j){
			int cur=frame[14-i+j][14-j];
			if(!incal){
				if(cur==0) // ---
					lsp++;
				else{
					incal=cur;// --*
					cnt=1;
					if(j==i)
						addcond(&s,incal,cnt,lsp,0,0,nextturn);
				}
			}else{ // incal!=0
				if(cur==0){ // -
					if(last==0){ // --
						int k;
						assert(msp!=0); //  --
						if(msp==cnt)
							msp=0;
						for(k=j;k<=i && !frame[14-i+k][14-k];++k)
							rsp++;
						addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						cnt=rsp=incal=0;
						lsp=2;
					}else{ // last!=0 ***-
						if(j==i){
							addcond(&s,incal,cnt,lsp,1,msp,nextturn);
							continue;
						}
						if(msp==0){ 
							msp=cnt;
							rsp=1;
						}else{ // msp!=0 && last!=0: *-**-
							if(frame[14-i+j+1][14-j-1]==0){ // *-**--
								addcond(&s,incal,cnt,lsp,2,msp,nextturn);
								cnt=incal=rsp=0;
								lsp=1;
							}else{ // *-**-?
								addcond(&s,incal,cnt,lsp,1,msp,nextturn);
								if(frame[14-i+j+1][14-j-1]==incal){ // *-**-*
									cnt=cnt-msp;
									lsp=1;
									msp=cnt;
									rsp=1;
								}else{ // *--**-x
									cnt=0;
									lsp=1;
									msp=0;
									rsp=0;
									incal=0;
								}
							}
						}
					}
				}else{	// incal!=0, cur!=0 : *x || * x ||  ** ||  *-*
					if(j==i){
						if(cur!=incal)
							addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						else
							addcond(&s,incal,cnt+1,lsp,0,msp,nextturn);
						continue;
					}
					if(incal==cur){
						cnt++;
						rsp=0;
					}else{
						addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						incal=cur;
						lsp=rsp;
						rsp=0;
						msp=0;
						cnt=1;
					}
					
				}
			}
			last=cur;
		}
	}




// half \

	for(i=4;i<15;++i){
		int incal=0,last=0; // 0,1,2
		int lsp=0,rsp=0,msp=0;
		int cnt=0;
		for(j=0;j<=i;++j){
			int cur=frame[j][14-i+j];
			if(!incal){
				if(cur==0) // ---
					lsp++;
				else{
					incal=cur;// --*
					cnt=1;
					if(j==i)
						addcond(&s,incal,cnt,lsp,0,0,nextturn);
				}
			}else{ // incal!=0
				if(cur==0){ // -
					if(last==0){ // --
						int k;
						assert(msp!=0); //  --
						if(msp==cnt)
							msp=0;
						for(k=j;k<=i && !frame[k][14-i+k];++k)
							rsp++;
						addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						cnt=rsp=incal=0;
						lsp=2;
					}else{ // last!=0 ***-
						if(j==i){
							addcond(&s,incal,cnt,lsp,1,msp,nextturn);
							continue;
						}
						if(msp==0){ 
							msp=cnt;
							rsp=1;
						}else{ // msp!=0 && last!=0: *-**-
							if(frame[j+1][14-i+j+1]==0){ // *-**--
								addcond(&s,incal,cnt,lsp,2,msp,nextturn);
								cnt=incal=rsp=0;
								lsp=1;
							}else{ // *-**-?
								addcond(&s,incal,cnt,lsp,1,msp,nextturn);
								if(frame[j+1][14-i+j+1]==incal){ // *-**-*
									cnt=cnt-msp;
									lsp=1;
									msp=cnt;
									rsp=1;
								}else{ // *--**-x
									cnt=0;
									lsp=1;
									msp=0;
									rsp=0;
									incal=0;
								}
							}
						}
					}
				}else{	// incal!=0, cur!=0 : *x || * x ||  ** ||  *-*
					if(j==i){
						if(cur!=incal)
							addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						else
							addcond(&s,incal,cnt+1,lsp,0,msp,nextturn);
						continue;
					}
					if(incal==cur){
						cnt++;
						rsp=0;
					}else{
						addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						incal=cur;
						lsp=rsp;
						rsp=0;
						msp=0;
						cnt=1;
					}
					
				}
			}
			last=cur;
		}
	}
// left part of \

	for(i=13;i>=4;--i){
		int incal=0,last=0; // 0,1,2
		int lsp=0,rsp=0,msp=0;
		int cnt=0;
		for(j=0;j<=i;++j){
			int cur=frame[14-i+j][j];
			if(!incal){
				if(cur==0) // ---
					lsp++;
				else{
					incal=cur;// --*
					cnt=1;
					if(j==i)
						addcond(&s,incal,cnt,lsp,0,0,nextturn);
				}
			}else{ // incal!=0
				if(cur==0){ // -
					if(last==0){ // --
						int k;
						assert(msp!=0); //  --
						if(msp==cnt)
							msp=0;
						for(k=j;k<=i && !frame[14-i+k][k];++k)
							rsp++;
						addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						cnt=rsp=incal=0;
						lsp=2;
					}else{ // last!=0 ***-
						if(j==i){
							addcond(&s,incal,cnt,lsp,1,msp,nextturn);
							continue;
						}
						if(msp==0){ 
							msp=cnt;
							rsp=1;
						}else{ // msp!=0 && last!=0: *-**-
							if(frame[14-i+j+1][j+1]==0){ // *-**--
								addcond(&s,incal,cnt,lsp,2,msp,nextturn);
								cnt=incal=rsp=0;
								lsp=1;
							}else{ // *-**-?
								addcond(&s,incal,cnt,lsp,1,msp,nextturn);
								if(frame[14-i+j+1][j+1]==incal){ // *-**-*
									cnt=cnt-msp;
									lsp=1;
									msp=cnt;
									rsp=1;
								}else{ // *--**-x
									cnt=0;
									lsp=1;
									msp=0;
									rsp=0;
									incal=0;
								}
							}
						}
					}
				}else{	// incal!=0, cur!=0 : *x || * x ||  ** ||  *-*
					if(j==i){
						if(cur!=incal)
							addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						else
							addcond(&s,incal,cnt+1,lsp,0,msp,nextturn);
						continue;
					}
					if(incal==cur){
						cnt++;
						rsp=0;
					}else{
						addcond(&s,incal,cnt,lsp,rsp,msp,nextturn);
						incal=cur;
						lsp=rsp;
						rsp=0;
						msp=0;
						cnt=1;
					}
					
				}
			}
			last=cur;
		}
	}


	return s;
}

int getpossible(struct STEP st[MAX_STEP]){
	int i,j;
	int n=0;
	if(curstep==0){
		st[0].x=7;
		st[0].y=7;
		st[0].bw=1;
		return 1;
	}
	for(i=0;i<15;++i)
		for(j=0;j<15;j++){
			if(frame[i][j]==0){
				int x,y;
				int find;
				for(x=-2,find=0;x<=2;++x){
					if(i+x<0 || i+x>14) continue;
					for(y=-2;y<=2;y++){
						if(j+y<0 || j+y>14 || (x==0 && y==0))
							continue;	
						if(frame[x+i][y+j]){
							find=1;
							st[n].x=i;
							st[n].y=j;
							st[n].bw=computerbw;//(curstep%2)?2:1;
							n++;
							break;
						}
					}
					if(find)
						break;
				}
			}
		}
	return n;
}


void applystep(struct STEP st){
	assert(frame[st.x][st.y]==0);
	frame[st.x][st.y]=st.bw;
}

void undostep(struct STEP st){
	frame[st.x][st.y]=0;
}

struct STEP algol2(){
	struct STEP ret[MAX_STEP];
	struct STEP steps[MAX_STEP];
	
}

struct STEP algol1(){
	int nret;
	struct STEP ret[MAX_STEP];
	struct STEP steps[MAX_STEP];
	int max=SCORE_INIT;
	int npos=getpossible(steps);
	int i;
	int samecnt=0;
	int next=computerbw==2?1:2;
	for(i=0;i<npos;i++){
		BWCORE sc;
		sc.bscore=sc.wscore=0;
		int score;
		applystep(steps[i]);
if(steps[i].x==6 && steps[i].y==5)
	printf("enter watch\n");
		sc=evaluate(next);
printf("%d,%d: b: %d, w %d\n",steps[i].x,steps[i].y,sc.bscore,sc.wscore);
		if(computerbw==2)
			score=sc.wscore-sc.bscore;
		else
			score=sc.bscore-sc.wscore;
		if(score>max){
			max=score;
			ret[0]=steps[i];
			samecnt=1;
		}else if(score==max){
			ret[samecnt++]=steps[i];
		}
		undostep(steps[i]);
	}
	npos=0;
	if(samecnt>1)
		npos=rand()%samecnt;
    applystep(ret[npos]);
	printf("%d,%d, evaluate %d\n",ret[npos].x,ret[npos].y,max);
	return ret[npos];
}

void init_robot(int bw,int level){
	srand(time(NULL));
	memset(frame,0,MAX_STEP);
	computerbw=bw;
	ailevel=0;//leve;
}

struct STEP get_step(){
	struct STEP st;
	if (ailevel==0)
		st=algol1();
	// else st=minmax(ailevel);
	steplog[curstep++]=st;
	return st;
}

void set_step(struct STEP st){
	frame[st.x][st.y]=st.bw;
	steplog[curstep++]=st;	
}

int retreat(){
	if (curstep>1)
	{
		frame[steplog[curstep-1].x][steplog[curstep-1].y]=0;
		frame[steplog[curstep-2].x][steplog[curstep-2].y]=0;
		curstep-=2;
	}
	return curstep;
}


