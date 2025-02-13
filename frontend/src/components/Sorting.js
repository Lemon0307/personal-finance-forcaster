const quickSort = (array, key, order) => {
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

    return [...quickSort(left, key, order), ...equal, ...quickSort(right, key, order)];
}

export default quickSort;