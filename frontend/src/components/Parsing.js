export const parseCSVToJSON = (file) => {
    // convert csv file into array of rows
    const lines = file.trim().split("\n")
    const transactions = []

    for (let i = 1; i < lines.length; i++) {
        // extract row column
        const row = lines[i].trim().replace(/^,/, "")
        const col = row.split(",").map(i => i.trim())
        // turn one row into a transaction object
        const transaction = {
            date: parseDate(col[0]),
            type: col[1].toLowerCase(),
            name: col[2],
            amount: parseFloat(col[3])
        }
        // push transaction object onto the array
        transactions.push(transaction)
    }

    return transactions
}

// DD/MM/YYYY to YYYY-MM-DD

const parseDate = (date) => {
    const temp = date.split("/")
    temp.reverse()
    return temp.join("-")
}