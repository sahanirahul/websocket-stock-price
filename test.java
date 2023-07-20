// Reverse a 32-bit signed integer and make sure the result is within the range of a 32-bit signed integer. 
// If the result overflows, return 0. Note that we should not use string variable and any variable that is more than 32-bit. 

// TEST 1 
// Input: 123 
// Output: 321
// TEST 2
// Input: -123 
// Output: -321

// 00001101 -> 13
// 00001011 ->
// 00000001 -> 00000010
// 

// 123 % 10 
// 123 / 10 -> 12

// 2^31 - 1
// -2^31
// 300 + 20 + 1 = 321

public int reverse(int num){
    boolean neg = false;
    if(num < 0){
        neg = true;
        num = -1*num;
    }

    int digits = 0;
    int tmp = num;
    while(tmp > 0){
        tmp = tmp / 10;
        digits++;
    }
    int ans = 0;
    while(num > 0){
        int d = num % 10;
        ans = ans + d*(int) Math.pow(10,digits - 1);
        if(ans < 0){
            return 0;
        }
        num = num/10;
        digits--;
    }
    if(neg) return -1*ans;
    return ans;
}

//Input: s1 = "pale", s2 = "ple"
//Input: s1 = "paled", s2 = "paler"
//Output: True


public int check(String s1,String s2){
    int m = s2.length();
    int n = s1.length();
    i
    int i = 0;
    while( i < s1.length){
        if(s1.charAt(i) != s2.charAt(i)){
            //insert
            if(s1.substring(i + 1).equals(s2.substring(i))) return true;
            if(s1.substring(i).equals(s2.substring(i + 1))) return true;
            if(s1.substring(i + 1).equals(s2.substring(i + 1))) return true;
            return false;
        }
        i++;
    }
    return true;
}