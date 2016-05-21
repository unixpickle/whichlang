# whichlang

This is a suite of Machine Learning tools for identifying the language in which a piece of code is written. It could potentially be used for text editors, code hosting websites, and much more.

A seasoned programmer could quickly tell you that this program is written in C:

```c
#include <stdio.h>

int main(int argc, const char ** argv) {
  printf("Hello, world!");
}
```

The goal of `whichlang` is to teach a program to do the same. By showing a Machine Learning algorithm a ton of code, it can *learn* to identify programming languages itself.

# Usage

There are four steps to using whichlang:

 * Configure Go and download whichlang.
 * Fetch code samples from Github or some other source.
 * Train a classifier with the code samples.
 * Use the whichlang API or server with the classifier you trained.

## Configuring Go and whichlang

First, follow the instructions on [this page](https://golang.org/doc/install) to setup Go. Once Go is setup and you have a `GOPATH` configured, run this set of commands:

```
$ go get github.com/unixpickle/whichlang
$ cd $GOPATH/src/github.com/unixpickle/whichlang
```

Now you have downloaded `whichlang` and are sitting in its root source folder.

## Fetching samples

To fetch samples from Github, you must have a Github account (having more than one Github account may be beneficial, as well). You should decide how many samples you want for each programming language. I have found that 180 is more than enough.

You can fetch samples and save them to a directory as follows:

```
$ mkdir /path/to/samples
$ go run cmd/fetchlang/*.go /path/to/samples 180
```

In the above example, I specified 180 samples per language. This will prompt you for your Github credentials (to get around strict API rate limits). If you specify a large number of samples (where 180 counts as a large number), you may hit Github's API rate limits several times during the fetching process. If this occurs, you will want to delete the partially-downloaded source directories (they will be subdirectories of your sample directory, and will contain less than 180 samples), then wait an hour before re-running `fetchlang`. The `fetchlang` sub-command will automatically skip any source directories that are already present, making it relatively easy to resume paused or rate-limited downloads.

## Training a classifier

With whichlang, you can train a number of different kinds of classifiers on your data. Currently, you can use the following classifiers:

 * [ID3](https://en.wikipedia.org/wiki/ID3)
 * [K-nearest neighbors](https://en.wikipedia.org/wiki/K-nearest_neighbors_algorithm)
 * [Artificial Neural Networks](https://en.wikipedia.org/wiki/Artificial_neural_network)
 * [Support Vector Machines](https://en.wikipedia.org/wiki/Support_vector_machine)

Out of these algorithms, I have found that Support Vector Machines are the simplest to train and work very well. Artificial Neural Networks are a close second, but they have more hyper-parameters and are thus harder to tune well. In this document, I will describe how to train both of these classifiers, leaving out ID3 and K-nearest neighbors.

### Choosing the "ubiquity"

For any classifier you use, you must choose a "ubiquity" value. Since whichlang works by extracting keywords from source files, it is important to discern potentially important keywords from file-specific keywords like variable names or embedded strings. To do this, keywords which appear in less than `N` files are ignored during training and classification, where `N` is the "ubiquity". I have found that a ubiquity of 10-20 works when you have roughly 100 source files.

### Support Vector Machines

The most basic way to train a Support Vector Machine is to allow whichlang to select all the hyper-parameters for you. Note, however, that this option is *very* slow, so you may want to keep reading.

```
$ go run cmd/trainer/*.go svm 15 /path/to/samples /path/to/classifier.json
```

In the above command, I specified a ubiquity of 15 files. This command will go through many different possible SVM configurations, choosing the one which performs the best on new samples (as measured via cross-validation). Since this command has to try many possible configurations, it will take a long time to run. I have already gone through the trouble of finding ideal parameters, and I will now share my results.

I have found that a linear SVM works fairly well for programming language classification. In particular, I've gotten a linear SVM to achieve a 93% success rate on new samples, and most of those mistakes were reasonable (e.g. mistaking C for C++, or mistaking Ruby for CoffeeScript). To train a linear SVM, you can set the `SVM_KERNEL` environment variable before running the `trainer` sub-command:

```
$ export SVM_KERNEL=linear
```

If you want verbose output during training, you can specify another environment variable:

```
$ export SVM_VERBOSE=1
```

For other SVM environment variables (e.g. for other kernels) you can checkout [this list](https://godoc.org/github.com/unixpickle/whichlang/svm#pkg-constants).

### Artificial Neural Networks

While whichlang does allow you to train ANNs without specifying any hyper-parameters (via grid search), doing so will take a tremendous amount of time. It is highly recommended that you manually specify the parameters for your neural network. I will give one example of training an ANN, but it is up to you to tweak these parameters:

```
$ export NEURALNET_VERBOSE=1
$ export NEURALNET_VERBOSE_STEPS=1
$ export NEURALNET_STEP_SIZE=0.01
$ export NEURALNET_MAX_ITERS=100
$ export NEURALNET_HIDDEN_SIZE=150
$ go run cmd/trainer/*.go neuralnet 15 /path/to/samples /path/to/classifier.json
```

## Using a classifier

Using a classifier is as simple as loading in a file. You can checkout the [classify command](https://github.com/unixpickle/whichlang/blob/master/cmd/classify/main.go) for a very simple (15-line) example.
