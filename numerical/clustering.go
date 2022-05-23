package numerical

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/edotau/goFish/csv"
	"github.com/edotau/goFish/stats"
)

type KMeans struct {
	// maxIterations is the number of iterations the learning will be cut off at in a non-online setting.
	maxIterations int
	// alpha is only used in the online setting of the algorithm
	alpha float64
	// trainingSet and guesses are the/ 'x', and 'y' of the data, expressed as vectors, that the model can optimize from.
	// Note that because K-Means is an unsupervised algorithm, the 'guesses' parameter is set while learning.
	// If you want to use the training not only to predict but just cluster an existing dataset, this storage will let the user export the predictions [][]float64{guesses[i]} == Predict(trainingSet[i])
	trainingSet [][]float64
	guesses     []int
	Centroids   [][]float64 `json:"centroids"`
	Output      io.Writer
}

/*
KNN implements the KNN algorithm
for classification, where an input
is classified by finding the K
nearest (by some distance metric)
data points, and taking a vote
based on those.
https://en.wikipedia.org/wiki/K-nearest_neighbors_algorithm
Example K-Nearest-Neighbors Model Usage:
	// initialize data!
	twoClusters := [][]float64{}
	twoClustersY := []float64{}
	for i := -10.0; i < -3; i += 0.1 {
		for j := -10.0; j < 10; j += 0.1 {
			twoClusters = append(twoClusters, []float64{i, j})
			twoClustersY = append(twoClustersY, 0.0)
		}
	}
	for i := 3.0; i < 10; i += 0.1 {
		for j := -10.0; j < 10; j += 0.1 {
			twoClusters = append(twoClusters, []float64{i, j})
			twoClustersY = append(twoClustersY, 1.0)
		}
	}
	// create the model using 3 nearest neighbors
	// for prediction, using the Euclidean Distance
	// as the distance metric.
	model := NewKNN(3, twoClusters, twoClustersY, EuclideanDistance)
	// make predictions like usual
	guess, err := model.Predict([]float64{-10,1})
	if err != nil {
		panic("THERE WAS AN ERROR")
	}
	// update the K used (use 10 neighbors now)
	model.K = 10
*/
type KNN struct {
	// Distance holds the distance
	// measure for the KNN algorithm,
	// which is just a function that
	// maps 2 float64 vectors to a
	// float64
	Distance DistanceMeasure

	// K is the number of nearest
	// neighbors to classify based
	// on in the KNN prediction
	// algorithm
	K int

	// trainingSet holds all training
	// examples, while expectedResults
	// holds the associated class of the
	// corresponding example.
	trainingSet     [][]float64
	expectedResults []float64
}

// OnlineParams is used to pass optional parameters in to creating a new K-Means model if you want to learn using the online version of the model
type OnlineParams struct {
	Alpha    float64
	Features int
}

// NewKMeans returns a pointer to the k-means
// model, which clusters given inputs in an
// unsupervised manner. The algorithm only has
// one optimization method (unless learning with
// the online variant which is more of a generalization
// than the same algorithm) so you aren't allowed
// to pass one in as an option.
//
// n is an optional parameter which (if given) assigns
// the length of the input vector.
func NewKMeans(k, maxIterations int, trainingSet [][]float64, params ...OnlineParams) *KMeans {
	var features int
	if len(params) != 0 {
		features = params[0].Features
	} else if len(trainingSet) != 0 {
		features = len(trainingSet[0])
	}

	alpha := 0.5
	if len(params) != 0 {
		alpha = params[0].Alpha
	}

	// start all guesses with the zero vector.
	// they will be changed during learning
	var guesses []int
	guesses = make([]int, len(trainingSet))

	rand.Seed(time.Now().UTC().Unix())
	centroids := make([][]float64, k)
	for i := range centroids {
		centroids[i] = make([]float64, features)
		for j := range centroids[i] {
			centroids[i][j] = 10 * (rand.Float64() - 0.5)
		}
	}

	return &KMeans{
		maxIterations: maxIterations,

		alpha: alpha,

		trainingSet: trainingSet,
		guesses:     guesses,

		Centroids: centroids,
		Output:    os.Stdout,
	}
}

// nn represents an encapsulation
// of the Nearest Neighbor data for
// each datapoint to facilitate easy
// sorting
type nn struct {
	X []float64
	Y float64

	Distance float64
}

// NewKNN returns a pointer to the k-means
// model, which clusters given inputs in an
// unsupervised manner. The algorithm only has
// one optimization method (unless learning with
// the online variant which is more of a generalization
// than the same algorithm) so you aren't allowed
// to pass one in as an option.
//
// n is an optional parameter which (if given) assigns
// the length of the input vector.
func NewKNN(k int, trainingSet [][]float64, expectedResults []float64, distanceMeasure DistanceMeasure) *KNN {
	return &KNN{
		Distance:        distanceMeasure,
		K:               k,
		trainingSet:     trainingSet,
		expectedResults: expectedResults,
	}
}

// LearningRate returns the learning rate α for gradient
// descent to optimize the model. Could vary as a function
// of something else later, potentially.
func (k *KMeans) LearningRate() float64 {
	return k.alpha
}

// Examples returns the number of training examples (m)
// that the model currently is training from.
func (k *KMeans) Examples() int {
	return len(k.trainingSet)
}

// MaxIterations returns the number of maximum iterations
// the model will go through in GradientAscent, in the
// worst case
func (k *KMeans) MaxIterations() int {
	return k.maxIterations
}

// Predict takes in a variable x (an array of floats,) and
// finds the value of the hypothesis function given the
// current parameter vector θ
//
// if normalize is given as true, then the input will
// first be normalized to unit length. Only use this if
// you trained off of normalized inputs and are feeding
// an un-normalized input
func (k *KMeans) Predict(x []float64, normalize ...bool) ([]float64, error) {
	if len(x) != len(k.Centroids[0]) {
		return nil, fmt.Errorf("Error: Centroid vector should be the same length as input vector!\n\tLength of x given: %v\n\tLength of centroid: %v\n", len(x), len(k.Centroids[0]))
	}

	if len(normalize) != 0 && normalize[0] {
		stats.NormalizePoint(x)
	}

	var guess int
	minDiff := diff(x, k.Centroids[0])
	for j := 1; j < len(k.Centroids); j++ {
		difference := diff(x, k.Centroids[j])
		if difference < minDiff {
			minDiff = difference
			guess = j
		}
	}

	return []float64{float64(guess)}, nil
}

// Learn takes the struct's dataset and expected results and runs
// batch gradient descent on them, optimizing theta so you can
// predict based on those results
//
// This batch version of the model uses the k-means++
// instantiation method to generate a consistantly better
// model than regular, randomized instantiation of
// centroids.
// Paper: http://ilpubs.stanford.edu:8090/778/1/2006-13.pdf
func (k *KMeans) Learn() error {
	if k.trainingSet == nil {
		err := fmt.Errorf("ERROR: Attempting to learn with no training examples!\n")
		fmt.Fprintf(k.Output, err.Error())
		return err
	}

	examples := len(k.trainingSet)
	if examples == 0 || len(k.trainingSet[0]) == 0 {
		err := fmt.Errorf("ERROR: Attempting to learn with no training examples!\n")
		fmt.Fprintf(k.Output, err.Error())
		return err
	}

	centroids := len(k.Centroids)
	features := len(k.trainingSet[0])

	fmt.Fprintf(k.Output, "Training:\n\tModel: K-Means++ Classification\n\tTraining Examples: %v\n\tFeatures: %v\n\tClasses: %v\n...\n\n", examples, features, centroids)

	// instantiate the centroids using k-means++
	k.Centroids[0] = k.trainingSet[rand.Intn(len(k.trainingSet))]

	distances := make([]float64, len(k.trainingSet))
	for i := 1; i < len(k.Centroids); i++ {
		var sum float64
		for j, x := range k.trainingSet {
			minDiff := diff(x, k.Centroids[0])
			for l := 1; l < i; l++ {
				difference := diff(x, k.Centroids[l])
				if difference < minDiff {
					minDiff = difference
				}
			}

			distances[j] = minDiff * minDiff
			sum += distances[j]
		}

		target := rand.Float64() * sum
		j := 0
		for sum = distances[0]; sum < target; sum += distances[j] {
			j++
		}
		k.Centroids[i] = k.trainingSet[j]

	}

	iter := 0
	for ; iter < k.maxIterations; iter++ {

		// set new guesses
		//
		// store counts when assigning classes
		// so you won't have to sum them again later
		classTotal := make([][]float64, centroids)
		classCount := make([]int64, centroids)

		for j := range k.Centroids {
			classTotal[j] = make([]float64, features)
		}

		for i, x := range k.trainingSet {
			k.guesses[i] = 0
			minDiff := diff(x, k.Centroids[0])
			for j := 1; j < len(k.Centroids); j++ {
				difference := diff(x, k.Centroids[j])
				if difference < minDiff {
					minDiff = difference
					k.guesses[i] = j
				}
			}

			classCount[k.guesses[i]]++
			for j := range x {
				classTotal[k.guesses[i]][j] += x[j]
			}
		}

		newCentroids := append([][]float64{}, k.Centroids...)
		for j := range k.Centroids {
			// if no objects are in the same class,
			// reinitialize it to a random vector
			if classCount[j] == 0 {
				for l := range k.Centroids[j] {
					k.Centroids[j][l] = 10 * (rand.Float64() - 0.5)
				}
				continue
			}

			for l := range k.Centroids[j] {
				k.Centroids[j][l] = classTotal[j][l] / float64(classCount[j])
			}
		}

		// only update if something was deleted
		if len(newCentroids) != len(k.Centroids) {
			k.Centroids = newCentroids
		}
	}

	fmt.Fprintf(k.Output, "Training Completed in %v iterations.\n%v\n", iter, k)

	return nil
}

// diff returns the square magnitude of the
// vector subtraction between u and v. This
// is also known as the Squared Euclidean
// Distance:
//
// diff(u, v) == |u - v|^2
//
// **NOTE** The function assumes that u and
// v are the same dimension to avoid constant
// checking from within algorithms.
func diff(u, v []float64) float64 {
	var sum float64
	for i := range u {
		sum += (u[i] - v[i]) * (u[i] - v[i])
	}

	return sum
}

// Examples returns the number of training examples (m)
// that the model currently is holding
func (k *KNN) Examples() int {
	return len(k.trainingSet)
}

// insertSorted takes a array v, and inserts u into
// the list in the position such that the list is
// sorted inversely. The function will not change the length
// of v, though, such that if u would appear last
// in the combined sorted list it would just be omitted.
//
// if the length of V is less than K, then u is inserted
// without deleting the last element
//
// Assumes v has been sorted. Uses binary search.
func insertSorted(u nn, v []nn, K int) []nn {
	low := 0
	high := len(v) - 1
	for low <= high {
		mid := (low + high) / 2
		if u.Distance < v[mid].Distance {
			high = mid - 1
		} else if u.Distance >= v[mid].Distance {
			low = mid + 1
		}
	}

	if low >= len(v) && len(v) >= K {
		return v
	}

	sorted := append(v[:low], append([]nn{u}, v[low:]...)...)

	if len(v) < K {
		return sorted
	}
	return sorted[:len(v)]
}

// round rounds a float64
func round(a float64) float64 {
	if a < 0 {
		return math.Ceil(a - 0.5)
	}
	return math.Floor(a + 0.5)
}

// Predict takes in a variable x (an array of floats,) and
// finds the value of the hypothesis function given the
// current parameter vector θ
//
// if normalize is given as true, then the input will
// first be normalized to unit length. Only use this if
// you trained off of normalized inputs and are feeding
// an un-normalized input
func (k *KNN) Predict(x []float64, normalize ...bool) ([]float64, error) {
	if k.K > len(k.trainingSet) {
		return nil, fmt.Errorf("Given K (%v) is greater than the length of the training set", k.K)
	}
	if len(x) != len(k.trainingSet[0]) {
		return nil, fmt.Errorf("Given x (len %v) does not match dimensions of training set", len(x))
	}

	if len(normalize) != 0 && normalize[0] {
		stats.NormalizePoint(x)
	}

	// initialize the neighbors as an empty
	// slice of Neighbors. insertSorted will
	// take care of capping the neighbors at
	// K.
	neighbors := []nn{}

	// calculate nearest neighbors
	for i := range k.trainingSet {
		dist := k.Distance(x, k.trainingSet[i])
		neighbors = insertSorted(nn{
			X: k.trainingSet[i],
			Y: k.expectedResults[i],

			Distance: dist,
		}, neighbors, k.K)
	}

	// take weighted vote
	sum := 0.0
	for i := range neighbors {
		sum += neighbors[i].Y
	}

	return []float64{round(sum / float64(k.K))}, nil
}

type DistanceMeasure func([]float64, []float64) float64

// EuclideanDistance returns the distance
// betweek two float64 vectors. NOTE that
// this function does not check that the
// vectors are different lengths (to improve
// computation speed in, say, KNN.) Make
// sure you pass in same-length vectors.
func EuclideanDistance(u []float64, v []float64) float64 {
	var sum float64
	for i := range u {
		sum += (u[i] - v[i]) * (u[i] - v[i])
	}
	return math.Sqrt(sum)
}

func (k *KMeans) OnlineLearn(errors chan error, dataset chan Datapoint, onUpdate func([][]float64), normalize ...bool) {
	if errors == nil {
		errors = make(chan error)
	}
	if dataset == nil {
		errors <- fmt.Errorf("ERROR: Attempting to learn with a nil data stream!\n")
		close(errors)
		return
	}

	centroids := len(k.Centroids)
	features := len(k.Centroids[0])

	fmt.Fprintf(k.Output, "Training:\n\tModel: Online K-Means Classification\n\tFeatures: %v\n\tClasses: %v\n...\n\n", features, centroids)

	var point Datapoint
	var more bool

	oneMinusAlpha := 1.0 - k.alpha

	for {
		point, more = <-dataset

		if more {
			if len(point.X) != features {
				errors <- fmt.Errorf("ERROR: point.X must have the same dimensions as clusters (len %v). Point: %v", len(point.X), point)
			}
			minDiff := diff(point.X, k.Centroids[0])
			c := 0
			for j := 1; j < len(k.Centroids); j++ {
				difference := diff(point.X, k.Centroids[j])
				if difference < minDiff {
					minDiff = difference
					c = j
				}
			}

			for i := range k.Centroids[c] {
				k.Centroids[c][i] = k.alpha*point.X[i] + oneMinusAlpha*k.Centroids[c][i]
			}

			go onUpdate([][]float64{{float64(c)}, k.Centroids[c]})

		} else {
			fmt.Fprintf(k.Output, "Training Completed.\n%v\n\n", k)
			close(errors)
			return
		}
	}
}

// String implements the fmt interface for clean printing. Here
// we're using it to print the model as the equation h(θ)=...
// where h is the k-means hypothesis model
func (k *KMeans) String() string {
	return fmt.Sprintf("h(θ,x) = argmin_j | x[i] - μ[j] |^2\n\tμ = %v", k.Centroids)
}

// Guesses returns the hidden parameter for the
// unsupervised classification assigned during
// learning.
//
//    model.Guesses[i] = E[k.trainingSet[i]]
func (k *KMeans) Guesses() []int {
	return k.guesses
}

// Distortion returns the distortion of the clustering
// currently given by the k-means model. This is the
// function the learning algorithm tries to minimize.
//
// Distorition() = Σ |x[i] - μ[c[i]]|^2
// over all training examples
func (k *KMeans) Distortion() float64 {
	var sum float64
	for i := range k.trainingSet {
		sum += diff(k.trainingSet[i], k.Centroids[int(k.guesses[i])])
	}

	return sum
}

// SaveClusteredData takes operates on a k-means
// model, concatenating the given dataset with the
// assigned class from clustering and saving it to
// file.
//
// Basically just a wrapper for the SaveDataToCSV
// with the K-Means data.
func (k *KMeans) SaveClusteredData(filepath string) {
	floatGuesses := []float64{}
	for _, val := range k.guesses {
		floatGuesses = append(floatGuesses, float64(val))
	}

	csv.Write(filepath, k.trainingSet, floatGuesses, true)
}

// PersistToFile takes in an absolute filepath and saves the
// centroid vector to the file, which can be restored later.
// The function will take paths from the current directory, but
// functions
//
// The data is stored as JSON because it's one of the most
// efficient storage method (you only need one comma extra
// per feature + two brackets, total!) And it's extendable.
func (k *KMeans) PersistToFile(path string) error {
	if path == "" {
		return fmt.Errorf("ERROR: you just tried to persist your model to a file with no path!! That's a no-no. Try it with a valid filepath")
	}

	bytes, err := json.Marshal(k.Centroids)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, bytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// RestoreFromFile takes in a path to a centroid vector
// and assigns the model it's operating on's parameter vector
// to that.
//
// The path must ba an absolute path or a path from the current
// directory
//
// This would be useful in persisting data between running
// a model on data, or for graphing a dataset with a fit in
// another framework like Julia/Gadfly.
func (k *KMeans) RestoreFromFile(path string) error {
	if path == "" {
		return fmt.Errorf("ERROR: you just tried to restore your model from a file with no path! That's a no-no. Try it with a valid filepath")
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &k.Centroids)
	if err != nil {
		return err
	}

	return nil
}

// Datapoint is used in some models where it is cleaner
// to pass data as a struct rather than just as 1D and
// 2D arrays like Generalized Linear Models are doing,
// for example. X corresponds to the inputs and Y
// corresponds to the result of the hypothesis.
//
// This is used with the Perceptron, for example, so
// data can be easily passed in channels while staying
// encapsulated well.
type Datapoint struct {
	X []float64 `json:"x"`
	Y []float64 `json:"y"`
}
