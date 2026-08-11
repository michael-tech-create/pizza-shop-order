import math

def temperature(n):
    if n == 0:
        return 0
    
    # 1. Start with None so the first number (-5) can be captured properly
    close_temp = None

    for i in input().split():
        t = int(i)

        if close_temp is None:
            close_temp = t
            continue
            
        current_d1 = int(math.fabs(t))
        closest_dist = int(math.fabs(close_temp))

        if current_d1 < closest_dist:
            close_temp = t
        elif current_d1 == closest_dist:
            if t > close_temp:
                close_temp = t

    return close_temp

def main():
    n = temperature(int(input("Enter a number: ")))
    print(n)

if __name__ == "__main__":
    main()
