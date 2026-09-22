/* 这份 C 代码就是本题的题面，不参与编译。
 * 请把它按 README 里的十步清单改写成 translate/translate.go。 */

#include <stdio.h>

int square(int n) {
    return n * n;
}

int main(void) {
    int total = 0;

    for (int i = 1; i <= 5; i++) {
        if (i % 2 == 0) {
            printf("%d even %d\n", i, square(i));
        } else {
            printf("%d odd %d\n", i, square(i));
        }
        total = total + square(i);
    }

    printf("total %d\n", total);
    return 0;
}
