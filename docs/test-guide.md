# Tests and configuration

* The tests (including the score/setting files) must be uploaded in `.zip` file not exceeding 200MB
* For every test id there must be an input(.in) and output(.out/ok/sol) file

## Naming format
* `{test_id}-{suffix}.{in|out|ok|sol}`
* `{suffix}.{test_id}.{in|out|ok|sol}`
* `{test_id}.{in|out|ok|sol}`

## Test Scores File
* Alongside the tests you can upload a `.txt` file specifying the individual score for each test. If left out, the score is distributed evenly across all tests
* Each line corresponds to a test and should have the following format: {test_id} {test_score}

## Problem Settings File
* Alongside the tests you can upload a `.properties` file to set problem parameters

```ini
# Setup subtasks
groups = 0,1-5,6-10,11-18,19-30
weights = 0,10,10,25,55
# Include subtasks in other subtasks
# For example Subtask 2 includes 1, Subtask 3 includes 1 and 2, and Subtask 5 includes 1,2,3,4
dependencies = ,1,1;2,,1;2;3;4
time = 1.200
memory = 512
```

# Checker

* Must be named `checker_legacy.cpp` or `checker.cpp` and added as a private attachement.

## Legacy Format

```cpp
#include <bits/stdc++.h>

using namespace std;

int score = 0;

void result(const char* msg, int pts) {
    // "{score} {message}" to stdout
    // score is an integer between 0 and 100
	printf("%d %s", pts, msg);
	exit(0);
}

void Assert(bool cond, string str) {
	if (!cond)
		result(str.c_str(), 0);
}
void Success(int pts, string str) {
	result(str.c_str(), pts);
}

ifstream out, ok, in;

int main(int argc, char* argv[]) {
	out.open(argv[1]); // user output
	ok.open(argv[2]); // correct output
	in.open(argv[3]); // test input; assumed to be valid
    int task;
	in >> task;
    if(task == 1){
        int ansok, ansout;
        ok >> ansok;
        Assert(!!(out >> ansout), "Wrong format"); // try to read contestant's answer
        Assert(ansout == ansok, "Wrong answer"); // check if values match
        score += 100;
    }
    else {
        int nrok, nrout;
        // not giving partial points if the first part of the answer is wrong, even if the second one might be right
        Assert(!!(out >> nrout), "Wrong answer. Incorrect number");
        score += 50;
        string message = "Correct number. ";
        double valok, valout;
        ok >> valok;
        // giving partial points even if the second part of answer is wrong
        if(!!(out >> valout) && abs(valok - valout) <= 1e-5)
            score += 50, message += "Correct value";
        else
            message += "Wrong value";
        Success(score, message);
    }
    Success(score, "Correct answer");
}
```

## Standard Format

```cpp
#include <bits/stdc++.h>

using namespace std;

float score = 0;

void result(const char* msg, float pts) {
    //"{score}" to stdout
    //"{message}" to stderr
    // score is a float between 0 and 1
	printf("%f", pts);
    fprintf(stderr, "%s", msg);
	exit(0);
}

void Assert(bool cond, string str) {
	if (!cond)
		result(str.c_str(), 0);
}
void Success(float pts, string str) {
	result(str.c_str(), pts);
}

ifstream out, ok, in;

int main(int argc, char* argv[]) {
	in.open(argv[1]); // test input; assumed to be valid
	ok.open(argv[2]); // correct output
	out.open(argv[3]); // user output
    int task;
	in >> task;
    if(task == 1){
        int ansok, ansout;
        ok >> ansok;
        Assert(!!(out >> ansout), "Wrong format"); // try to read contestant's answer
        Assert(ansout == ansok, "Wrong answer"); // check if values match
        score += 0.5;
    }
    else {
        int nrok, nrout;
        // not giving partial points if the first part of the answer is wrong, even if the second one might be right
        Assert(!!(out >> nrout), "Wrong answer. Incorrect number");
        score += 0.5;
        string message = "Correct number. ";
        double valok, valout;
        ok >> valok;
        // giving partial points even if the second part of answer is wrong
        if(!!(out >> valout) && abs(valok - valout) <= 1e-5)
            score += 0.5, message += "Correct value";
        else
            message += "Wrong value";
        Success(score, message);
    }
    Success(score, "Correct answer");
}
```

# Grader

* Must be named `grader.cpp` and added as a private attachment.
* It must contain the main function.
* It can include other header files included as public (if they contain functions the user must implement/use) attachments.
* Optionally, a sample grader can be added as a public attachment to make local testing easier. Keep in mind it must be prefixed with an underscore.
* Custom graders are usually used alongside custom checkers.
* **If the input the grader reads must remain a secret, it is recommended you read it from a file and overwrite it before you call any of the user's functions.**
* **In some situations to preserve the integrity of the function calls, you might want to print a secret message (random string or hash of input) using the grader and check it before reading anything else using the custom checker.**

## Sample

`grader.cpp`
```cpp
#include <bits/stdc++.h>
#include "sum.h" // mandatory, as we're just using it without defining it; 
                 // if the user doesn't use functions defined in this file, we can get rid of the header file altogether and declare the functions here (that the user will implement) without defining them 

using namespace std;

int main() {
    int a, b;
    cin >> a >> b;
    cout << sum(a, b);
}
```

`sum.h`
```cpp
int sum(int a, int b);
```

`user`
```cpp
#include "sum.h" // not mandatory in this case, as we are just defining sum

int sum(int a, int b) {
    return a + b;
}
```