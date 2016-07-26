#include <stdio.h>
#include <pthread.h>
#include <sched.h>

pthread_mutex_t lock;
int ticket = 0;
int serving = 0;
int counter = 0;

int fetch_and_add()
{
	pthread_mutex_lock(&lock);
	int old;
	old = ticket;
	ticket++;
	pthread_mutex_unlock(&lock);
	return old;
}

int acquire() 
{
	int mytick = fetch_and_add(ticket);
	while (mytick != serving) {}	
	return mytick;
}

void release(int mytick)
{
	serving = mytick + 1;
}

void* small_func(void* thread_id)
{
	int k;
	for(k = 0; k < 1000000; k++)
	{
		int mytick = acquire();
		counter++;
		release(mytick);
	}
}

int main(int argc, char *argv[])
{
	int rc;
	pthread_t threads[2];

	long t;
	for(t = 0; t < 2; t++) {
		rc = pthread_create(&threads[t], NULL, small_func, (void*)t);
	}

	for (t = 0; t < 2; t++) {
		pthread_join(threads[t], NULL);
	}

	printf("counter is  %d\n", counter);
	pthread_exit(NULL);
}
