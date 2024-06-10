#include <list>
#include <future>
#include <iostream>
#include <algorithm>

template< typename T >
std::list<T> quick_sort(std::list<T> input) 
{
    if (input.empty()) return {}; // 特判

    std::list<T> ret;
    // 转移支点
    ret.splice(ret.end(), input, input.begin()); 
    auto piv = ret.front();

    // input重排序
    auto pivIt = std::partition(input.begin(), input.end(), [&piv](const auto& t) {
        return t < piv;
    });
    
    std::list<T> lowList;
    lowList.splice(lowList.end(), input, input.begin(), pivIt);

    auto newLow(quick_sort<T>(std::move(lowList)));
    auto newUp(quick_sort<T>(std::move(input)));

    ret.splice(ret.begin(), newLow);
    ret.splice(ret.end(), newUp);

    return ret;
}

template< typename T >
std::list<T> parallel_quick_sort(std::list<T> input) 
{
    if (input.empty()) return {}; // 特判

    std::list<T> ret;
    // 转移支点
    ret.splice(ret.end(), input, input.begin()); 
    auto piv = ret.front();

    // input重排序
    auto pivIt = std::partition(input.begin(), input.end(), [&piv](const auto& t) {
        return t < piv;
    });
    
    std::list<T> lowList;
    lowList.splice(lowList.end(), input, input.begin(), pivIt);

    // 该部分异步执行
    std::future<std::list<T>> newLow(std::async(quick_sort<T>, std::move(lowList)));
    
    auto newUp(quick_sort<T>(std::move(input)));

    ret.splice(ret.begin(), newLow.get());
    ret.splice(ret.end(), newUp);

    return ret;
}
std::ostream& operator<<(std::ostream& ostr, const std::list<int>& list)
{
    for (auto &i : list)
        ostr << ' ' << i;
 
    return ostr;
}


int main() {
    #if 0
    std::list ls{1, 2, 3, 4, 5};
    std::list<int> ls1;
    // 拼成环形链表了
    ls1.splice(ls1.end(), ls); // 拼接链表
    ls1.splice(ls1.end(), ls, ls.begin()); // 拼接结点
    std::cout << ls << "\n";
    std::cout << ls1 << "\n";
    #endif
    // std::cout << "-------------------\n";
    std::list l{1, 9, 2, 0, 3, 8, 4, 7, 5, 6};
    std::cout << parallel_quick_sort(l);
    return 0;
}