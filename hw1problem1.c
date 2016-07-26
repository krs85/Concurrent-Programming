#define _GNU_SOURCE
#include <stdio.h>
#include <pthread.h>
#include <sched.h>


#define NUM_THREADS	2

int counter = 0;
void* small_func(void* thread_id) 
{
	int k;
	for(k = 0; k < 1000000; k++) 
	{
		counter++;
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

	   
	long t;
	for(t = 0; t < NUM_THREADS; t++) { 
		rc = pthread_create(&threads[t], NULL, small_func, (void *)t);
	}

	for (t = 0; t < NUM_THREADS; t++) {
		pthread_join(threads[t], NULL);
	}
	printf("counter is  %d\n", counter);
	pthread_exit(NULL);
}

