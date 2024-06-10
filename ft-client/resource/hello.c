#include <stdio.h>
struct X {
    int a;
};

void f(int a) {
    return;
}
int main(){
    printf("aaa\n");
    int a = 1;
    f(a); 
    struct X x;
    struct X* p = &x;
    printf("%d", p->a);
    return 0;
}

