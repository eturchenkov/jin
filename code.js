function quickSort(arr) {
  if (arr.length <= 1) {
    return arr;
  }

  const pivot = arr[Math.floor(arr.length / 2)];

  const left = [];
  const right = [];
  const equal = [];

  for (const element of arr) {
    if (element < pivot) {
      left.push(element);
    } else if (element > pivot) {
      right.push(element);
    } else {
      equal.push(element);
    }
  }

  return [...quickSort(left), ...equal, ...quickSort(right)];
}

const unsortedArray = [3, 6, 8, 10, 1, 2, 1];
const sortedArray = quickSort(unsortedArray);
console.log(sortedArray);