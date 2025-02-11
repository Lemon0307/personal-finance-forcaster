const quickSort = (array, key) => {
    if (array.length <= 1) return array;

    const pivot = array[Math.floor(array.length / 2)];
    const left = [];
    const right = [];
    const equal = [];

    for (const item of array) {
        let comparison = 0
        if (typeof item[key] === "string" && typeof pivot[key] === "string") {
            comparison = item[key].localeCompare(pivot[key])
        } else {
            comparison = item[key] - pivot[key]
        }

        if (comparison < 0) {
            left.push(item);
        } else if (comparison > 0) {
            right.push(item);
        } else {
            equal.push(item);
        }
    }

    return [...quickSort(left, key), ...equal, ...quickSort(right, key)];
}

export default quickSort;