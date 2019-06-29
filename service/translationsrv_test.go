// 协慌网 - 专注编程问答汉化 https://routinepanic.com
// Copyright (C) 2018-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package service

import (
	"testing"
)

func TestTranslate(t *testing.T) {
	content := `
<p>To understand what <code translate="no">yield</code> does, you must understand what <em>generators</em> are. And before generators come <em>iterables</em>.</p>

<h2>Iterables</h2>

<p>When you create a list, you can read its items one by one. Reading its items one by one is called iteration:</p>

<pre translate="no"><code translate="no">&gt;&gt;&gt; mylist = [1, 2, 3]
&gt;&gt;&gt; for i in mylist:
...    print(i)
1
2
3
</code></pre>

<p><code translate="no">mylist</code> is an <em>iterable</em>. When you use a list comprehension, you create a list, and so an iterable:</p>

<pre translate="no"><code translate="no">&gt;&gt;&gt; mylist = [x*x for x in range(3)]
&gt;&gt;&gt; for i in mylist:
...    print(i)
0
1
4
</code></pre>

<p>Everything you can use &#34;<code translate="no">for... in...</code>&#34; on is an iterable; <code translate="no">lists</code>, <code translate="no">strings</code>, files...</p>

<p>These iterables are handy because you can read them as much as you wish, but you store all the values in memory and this is not always what you want when you have a lot of values.</p>

<h2>Generators</h2>

<p>Generators are iterators, a kind of iterable <strong>you can only iterate over once</strong>. Generators do not store all the values in memory, <strong>they generate the values on the fly</strong>:</p>

<pre translate="no"><code translate="no">&gt;&gt;&gt; mygenerator = (x*x for x in range(3))
&gt;&gt;&gt; for i in mygenerator:
...    print(i)
0
1
4
</code></pre>

<p>It is just the same except you used <code translate="no">()</code> instead of <code translate="no">[]</code>. BUT, you <strong>cannot</strong> perform <code translate="no">for i in mygenerator</code> a second time since generators can only be used once: they calculate 0, then forget about it and calculate 1, and end calculating 4, one by one.</p>

<h2>Yield</h2>

<p><code translate="no">yield</code> is a keyword that is used like <code translate="no">return</code>, except the function will return a generator.</p>

<pre translate="no"><code translate="no">&gt;&gt;&gt; def createGenerator():
...    mylist = range(3)
...    for i in mylist:
...        yield i*i
...
&gt;&gt;&gt; mygenerator = createGenerator() # create a generator
&gt;&gt;&gt; print(mygenerator) # mygenerator is an object!
&lt;generator object createGenerator at 0xb7555c34&gt;
&gt;&gt;&gt; for i in mygenerator:
...     print(i)
0
1
4
</code></pre>

<p>Here it&#39;s a useless example, but it&#39;s handy when you know your function will return a huge set of values that you will only need to read once.</p>

<p>To master <code translate="no">yield</code>, you must understand that <strong>when you call the function, the code you have written in the function body does not run.</strong> The function only returns the generator object, this is a bit tricky :-)</p>

<p>Then, your code will be run each time the <code translate="no">for</code> uses the generator.</p>

<p>Now the hard part:</p>

<p>The first time the <code translate="no">for</code> calls the generator object created from your function, it will run the code in your function from the beginning until it hits <code translate="no">yield</code>, then it&#39;ll return the first value of the loop. Then, each other call will run the loop you have written in the function one more time, and return the next value, until there is no value to return.</p>

<p>The generator is considered empty once the function runs, but does not hit <code translate="no">yield</code> anymore. It can be because the loop had come to an end, or because you do not satisfy an <code translate="no">&#34;if/else&#34;</code> anymore.</p>

<hr/>

<h2>Your code explained</h2>

<p>Generator:</p>

<pre translate="no"><code translate="no"># Here you create the method of the node object that will return the generator
def _get_child_candidates(self, distance, min_dist, max_dist):

    # Here is the code that will be called each time you use the generator object:

    # If there is still a child of the node object on its left
    # AND if distance is ok, return the next child
    if self._leftchild and distance - max_dist &lt; self._median:
        yield self._leftchild

    # If there is still a child of the node object on its right
    # AND if distance is ok, return the next child
    if self._rightchild and distance + max_dist &gt;= self._median:
        yield self._rightchild

    # If the function arrives here, the generator will be considered empty
    # there is no more than two values: the left and the right children
</code></pre>

<p>Caller:</p>

<pre translate="no"><code translate="no"># Create an empty list and a list with the current object reference
result, candidates = list(), [self]

# Loop on candidates (they contain only one element at the beginning)
while candidates:

    # Get the last candidate and remove it from the list
    node = candidates.pop()

    # Get the distance between obj and the candidate
    distance = node._get_dist(obj)

    # If distance is ok, then you can fill the result
    if distance &lt;= max_dist and distance &gt;= min_dist:
        result.extend(node._values)

    # Add the children of the candidate in the candidates list
    # so the loop will keep running until it will have looked
    # at all the children of the children of the children, etc. of the candidate
    candidates.extend(node._get_child_candidates(distance, min_dist, max_dist))

return result
</code></pre>

<p>This code contains several smart parts:</p>

<ul>
<li><p>The loop iterates on a list, but the list expands while the loop is being iterated :-) It&#39;s a concise way to go through all these nested data even if it&#39;s a bit dangerous since you can end up with an infinite loop. In this case, <code translate="no">candidates.extend(node._get_child_candidates(distance, min_dist, max_dist))</code> exhausts all the values of the generator, but <code translate="no">while</code> keeps creating new generator objects which will produce different values from the previous ones since it&#39;s not applied on the same node.</p></li>
<li><p>The <code translate="no">extend()</code> method is a list object method that expects an iterable and adds its values to the list.</p></li>
</ul>

<p>Usually we pass a list to it:</p>

<pre translate="no"><code translate="no">&gt;&gt;&gt; a = [1, 2]
&gt;&gt;&gt; b = [3, 4]
&gt;&gt;&gt; a.extend(b)
&gt;&gt;&gt; print(a)
[1, 2, 3, 4]
</code></pre>

<p>But in your code it gets a generator, which is good because:</p>

<ol>
<li>You don&#39;t need to read the values twice.</li>
<li>You may have a lot of children and you don&#39;t want them all stored in memory.</li>
</ol>

<p>And it works because Python does not care if the argument of a method is a list or not. Python expects iterables so it will work with strings, lists, tuples and generators! This is called duck typing and is one of the reason why Python is so cool. But this is another story, for another question...</p>

<p>You can stop here, or read a little bit to see an advanced use of a generator:</p>

<h2>Controlling a generator exhaustion</h2>

<pre translate="no"><code translate="no">&gt;&gt;&gt; class Bank(): # Let&#39;s create a bank, building ATMs
...    crisis = False
...    def create_atm(self):
...        while not self.crisis:
...            yield &#34;$100&#34;
&gt;&gt;&gt; hsbc = Bank() # When everything&#39;s ok the ATM gives you as much as you want
&gt;&gt;&gt; corner_street_atm = hsbc.create_atm()
&gt;&gt;&gt; print(corner_street_atm.next())
$100
&gt;&gt;&gt; print(corner_street_atm.next())
$100
&gt;&gt;&gt; print([corner_street_atm.next() for cash in range(5)])
[&#39;$100&#39;, &#39;$100&#39;, &#39;$100&#39;, &#39;$100&#39;, &#39;$100&#39;]
&gt;&gt;&gt; hsbc.crisis = True # Crisis is coming, no more money!
&gt;&gt;&gt; print(corner_street_atm.next())
&lt;type &#39;exceptions.StopIteration&#39;&gt;
&gt;&gt;&gt; wall_street_atm = hsbc.create_atm() # It&#39;s even true for new ATMs
&gt;&gt;&gt; print(wall_street_atm.next())
&lt;type &#39;exceptions.StopIteration&#39;&gt;
&gt;&gt;&gt; hsbc.crisis = False # The trouble is, even post-crisis the ATM remains empty
&gt;&gt;&gt; print(corner_street_atm.next())
&lt;type &#39;exceptions.StopIteration&#39;&gt;
&gt;&gt;&gt; brand_new_atm = hsbc.create_atm() # Build a new one to get back in business
&gt;&gt;&gt; for cash in brand_new_atm:
...    print cash
$100
$100
$100
$100
$100
$100
$100
$100
$100
...
</code></pre>

<p><strong>Note:</strong> For Python 3, use<code translate="no">print(corner_street_atm.__next__())</code> or <code translate="no">print(next(corner_street_atm))</code></p>

<p>It can be useful for various things like controlling access to a resource.</p>

<h2>Itertools, your best friend</h2>

<p>The itertools module contains special functions to manipulate iterables. Ever wish to duplicate a generator?
Chain two generators? Group values in a nested list with a one-liner? <code translate="no">Map / Zip</code> without creating another list?</p>

<p>Then just <code translate="no">import itertools</code>.</p>

<p>An example? Let&#39;s see the possible orders of arrival for a four-horse race:</p>

<pre translate="no"><code translate="no">&gt;&gt;&gt; horses = [1, 2, 3, 4]
&gt;&gt;&gt; races = itertools.permutations(horses)
&gt;&gt;&gt; print(races)
&lt;itertools.permutations object at 0xb754f1dc&gt;
&gt;&gt;&gt; print(list(itertools.permutations(horses)))
[(1, 2, 3, 4),
 (1, 2, 4, 3),
 (1, 3, 2, 4),
 (1, 3, 4, 2),
 (1, 4, 2, 3),
 (1, 4, 3, 2),
 (2, 1, 3, 4),
 (2, 1, 4, 3),
 (2, 3, 1, 4),
 (2, 3, 4, 1),
 (2, 4, 1, 3),
 (2, 4, 3, 1),
 (3, 1, 2, 4),
 (3, 1, 4, 2),
 (3, 2, 1, 4),
 (3, 2, 4, 1),
 (3, 4, 1, 2),
 (3, 4, 2, 1),
 (4, 1, 2, 3),
 (4, 1, 3, 2),
 (4, 2, 1, 3),
 (4, 2, 3, 1),
 (4, 3, 1, 2),
 (4, 3, 2, 1)]
</code></pre>

<h2>Understanding the inner mechanisms of iteration</h2>

<p>Iteration is a process implying iterables (implementing the <code translate="no">__iter__()</code> method) and iterators (implementing the <code translate="no">__next__()</code> method).
Iterables are any objects you can get an iterator from. Iterators are objects that let you iterate on iterables.</p>

<p>There is more about it in this article about <a href="http://effbot.org/zone/python-for-statement.htm" rel="noreferrer">how <code translate="no">for</code> loops work</a>.</p>
`

	text := Translation.Translate(content, "html")
	logger.Info(text)
}
