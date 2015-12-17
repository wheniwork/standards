At [When I Work](http://wheniwork.com/about), we do a fair amount of code reviews. We haven't always done them, but as we've grown they have become key to ensuring quality for the entire team. Here are many of the ways we approach them.

Briefly, a code review is a discussion between two or more developers about changes to the code to address an issue. Many articles talk about the benefits of code reviews, including knowledge sharing, code quality, and developer growth. There are significantly fewer that talk about what to look for in a review and how to actually execute a code review.

### What we look for during a review
#### Architecture / Design

* **[Single Responsibility Principle](http://en.wikipedia.org/wiki/Single_responsibility_principle):** The idea that a class should have one-and-only-one responsibility. Harder than one might expect. I usually apply this to methods too. If we have to use “and” to finish describing what a method is capable of doing, it might be at the wrong level of abstraction.

* **[Open/Closed Principle](http://en.wikipedia.org/wiki/Open/closed_principle):** If the language is object-oriented, are the objects open for extension but closed for modification? What happens if we need to add another one of `Class`?

* **Code Duplication:** Using the "three strikes" rule comes in handy. If code is copied once, it's usually okay, though still a bit painful. If it's copied again, it should be refactored so that the common functionality is split out.

* **[Squint-test Check](http://robertheaton.com/2014/06/20/code-review-without-your-eyes/):** This might seem funny, but can be very helpful. Are there any funny patterns/shapes that might indicate other problems in the code's structure?

* **Better Code:** If you're changing an area of the code that is messy, it's easy and tempting to add in a few lines and leave. We always encourage each other to go one step further and leaving the code nicer than it was found.

* **Potential Bugs:** Are there off-by-one errors? Will the loops terminate in the way we expect? Will they terminate at all?

* **Error Handling:** Are errors handled gracefully and explicitly where necessary? Have custom errors been added? If so, are they useful?

* **Efficiency:** If there's an algorithm in the code, is it using an efficient implementation? For example, iterating over a list of keys in a dictionary is an inefficient way to locate a desired value.

#### Style

* **Method Names:** Naming things is one of the hard problems in computer science. If a method is named `get_shift_list` and it is actually doing something completely different like sanitizing HTML from the input, then that's an inaccurate method name. And probably a misleading function.

* **Variable Names:** `$foo` or `$bar` are probably not useful names for data structures. `$e` is similarly not useful when compared to `$exception`. Be as verbose as you need (depending on the language). Expressive variable names make it easier to understand code when we have to revisit it later.

* **Function Size:** Our rule of thumb is that a function should be less than 50 or so lines. If we see a method above 75 or 100, we feel it’s best that it be cut into smaller pieces.

* **Class Size:** We think classes should be under about 500 lines total and ideally less than 300. It's likely that large classes can be split into separate classes or object, which makes it easier to understand what the class is supposed to do.

* **Docblocks:** For complex methods or those with longer lists of arguments, is there a docblock explaining what each of the arguments does, if it's not obvious?

* **Commented Out Code:** Seems like common sense, but don't keep around line of code that are commented out. We use git for a reason, right?

* **Method Arguments:** For the methods and functions, do they have a reasonable number arguments? If there are too many, it is probably a good sign that it could be grouped in a different way.

* **Readability:** Is the code easy to understand? Does its structure help to explain the purpose?

#### Testing

* **Test Coverage:** We like to see tests for new features. Are the tests thoughtful? Do they cover the failure conditions? Are they solid? Are they slow?

* **Meets Requirements:** Usually as part of the end of a review, we'll take a look at the requirements of the issue, task, or bug which the work was filed against. If it doesn't meet one of the criteria, it's better to bounce it back before it goes to QA.

### Review your own code first
Before submitting my code, I will often do a git add for the affected files / directories and then run a git diff --staged to examine the changes I have not yet committed. Usually I’m looking for things like:

* Did I leave a comment, `console` or TODO in?
* Do all the variable names make sense?
* ...and anything else that was brought up above.

We all want to make sure that the code would pass our own code review first before subjecting other people to it. It also stings less to get notes from yourself than from others. :stuck_out_tongue:

### How to handle code reviews
I find that the human parts of the code review offer as many challenges as reviewing the code. I’m still learning how to handle this part too. Here are some approaches that have worked for me when discussing code:

* **Ask questions:** How does this method work? If this requirement changes, what else would have to change? How could we make this more maintainable?

* **Compliment / reinforce good practices:** One of the most important parts of the code review is to reward developers for growth and effort. Few things feel better than getting praise from a peer. I try to offer as many positive comments as possible.

* **Discussions for more detailed points:** On occasion, a recommended architectural change might be large enough that it's easier to discuss it in real-time rather than in the comments.

* **Explain reasoning:** I find it’s best both to ask if there’s a better alternative and justify why I think it’s worth fixing. Sometimes it can feel like the changes suggested can seem nit-picky without context or explanation.

* **Make it about the code:** It’s easy to take notes from code reviews personally, especially if we take pride in our work. It’s best, I find, to make discussions about the code than about the developer. It lowers resistance and it’s not about the developer anyway, it’s about improving the quality of the code.

* **Suggest importance of fixes:** I tend to offer many suggestions, not all of which need to be acted upon. Clarifying if an item is important to fix before it can be considered done is useful both for the reviewer and the reviewee. It makes the results of a review clear and actionable.

### On mindset
As developers, we are responsible for making both working and maintainable code. It can be easy to defer the second part because of pressure to deliver working code. Refactoring does not change functionality by design, so don’t let suggested changes discourage you. Improving the maintainability of the code can be just as important as fixing the line of code that caused the bug.

In addition, please keep an open mind during code reviews. This is something I think everyone struggles with. I can get defensive in code reviews too, because it can feel personal when someone says code you wrote could be better.

If the reviewer makes a suggestion, and I don’t have a clear answer as to why the suggestion should not be implemented, I’ll usually make the change. If the reviewer is asking a question about a line of code, it may mean that it would confuse others in the future. In addition, making the changes can help reveal larger architectural issues or bugs.

### Addressing suggested changes
We typically leave comments on a per-line basis with some thinking behind them. Usually we will look at the review notes in Github and, at the same time, have the code pulled up to make the suggested changes. This way it's easier to remember and address all the items right away.
