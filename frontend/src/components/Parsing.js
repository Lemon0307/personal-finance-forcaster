export const parseCSVToJSON = (file) => {
    // convert csv file into array of rows
    const lines = file.trim().split("\n")
    const transactions = []

    const headers = lines[0].trim().split(",").map(h => h.toLowerCase())

    const required_columns = ["date", "type", "name", "amount"]

    const missing_columns = required_columns.filter(col => !headers.includes(col))
    if (missing_columns.length > 0) {
        alert("Please provide a CSV in the correct format: name, type, amount, date")
        return 1
    }

    const column_index = {
        date: headers.indexOf("date"),
        type: headers.indexOf("type"),
        name: headers.indexOf("name"),
        amount: headers.indexOf("amount")
    }

    for (let i = 1; i < lines.length; i++) {
        // extract row column
        const row = lines[i].trim().replace(/^,/, "")
        const col = row.split(",").map(i => i.trim())

        if (col.length < headers.length) continue

        // turn one row into a transaction object
        const transaction = {
            date: parseDate(col[column_index.date]),
            type: col[column_index.type].toLowerCase(),
            name: col[column_index.name],
            amount: parseFloat(col[column_index.amount])
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