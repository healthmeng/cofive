#ifndef __ROBOT_H__
#define __ROBOT_H__

#include <stdio.h>
#include <string.h>
#define MAX_STEP 15*15

struct STEP{
	int x,y;
	int bw;
//	int n;
};


void init_robot(int bw, int level);
struct STEP get_step();
void set_step(struct STEP st);
int retreat();

#endif
