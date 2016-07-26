#define _GNU_SOURCE
#include <stdio.h>
#include <pthread.h>
#include <sched.h>


#define NUM_THREADS	2

int counter = 0;

int wantp = 0;
int wantq = 0;
int turn = 1;

void* p_func(void* thread_id)
{
	int i = 0;
	for (i = 0; i < 1000000; i++) 
	{
		wantp = 1;
		while (wantq)
		{
			if (turn == 2) 
			{
				wantp = 0;
				/*await turn = 1*/
				while (turn != 1)
				{
					pthread_yield();
					//sched_yield();
					//sleep(0);
				}
				wantp = 1;
			}
		}
		counter++;
		turn = 2;
		wantp = 0;
	}
}

void* q_func(void* thread_id)
{
	int k = 0;
	for (k = 0; k < 1000000; k++) 
	{
		wantq = 1;
		while (wantp)
		{
			if (turn == 1) 
			{
				wantq = 0;
				/*await turn = 2*/
				while (turn != 2)
				{
					pthread_yield();
				}
				wantq = 1;
			}
		}
		counter++; //critical section
		turn = 1;
		wantq = 0;
	}
}

int main(int argc, char *argv[]) 
{
	pthread_attr_t myattr;
	cpu_set_t cpuset;

	pthread_attr_init(&myattr);
	CPU_ZERO(&cpuset);
	CPU_SET(0, &cpuset);
	pthread_attr_setaffinity_np(&myattr, sizeof(cpu_set_t), &cpuset);

	int rc;
	pthread_t threads[NUM_THREADS];

	rc = pthread_create(&threads[0], &myattr, p_func, (void*) 0);
	rc = pthread_create(&threads[1], &myattr, q_func, (void*) 1); 	
	pthread_join(threads[0], NULL);
	pthread_join(threads[1], NULL);
	//long t;
	//for(t = 0; t < NUM_THREADS; t++) { 
	//	rc = pthread_create(&threads[t], NULL, small_func, (void *)t);
	//}

	printf("counter is  %d\n", counter);
	pthread_exit(NULL);
}

