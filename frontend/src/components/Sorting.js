const quickSort = (data, key, order) => {
    // base case
    if (data.length <= 1) return data;

    // set pivot as the median of the data
    const pivot = data[Math.floor(data.length / 2)];
    const left = [];
    const right = [];
    const equal = [];

    for (const item of data) {
        let comparison = 0
        // check if item and pivot are strings
        if (typeof item[key] === "string" && typeof pivot[key] === "string") {
            // compare strings alphabetically
            comparison = item[key].localeCompare(pivot[key])
        } else {
            // compare numbers
            comparison = item[key] - pivot[key]
        }
        // put items onto left and right differently depending on
        // the order selected
        switch (order) {
            case "asc":
                if (comparison > 0) {
                    left.push(item);
                } else if (comparison < 0) {
                    right.push(item);
                } else {
                    equal.push(item);
                }                
                break
            case "desc":
                if (comparison < 0) {
                    left.push(item);
                } else if (comparison > 0) {
                    right.push(item);
                } else {
                    equal.push(item);
                }
                break
            default:
                break
        }
    }
    // recursively call quicksort on left and right
    return [...quickSort(left, key, order), ...equal, ...quickSort(right, key, order)];
}

export default quickSort;