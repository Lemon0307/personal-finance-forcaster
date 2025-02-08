export const parseCSVToJSON = (file) => {
    console.log(file)
    const lines = file.trim().split("\n")
    const transactions = []
    console.log(lines)

    for (let i = 1; i < lines.length; i++) {
        const row = lines[i].trim().replace(/^,/, "")
        const col = row.split(",").map(i => i.trim())
        const transaction = {
            date: parseDate(col[0]),
            type: col[1].toLowerCase(),
            name: col[2],
            amount: parseFloat(col[3])
        }
        transactions.push(transaction)
    }

    console.log(transactions)

    return transactions
}

// DD/MM/YYYY to YYYY-MM-DD

const parseDate = (date) => {
    const temp = date.split("/")
    temp.reverse()
    return temp.join("-")
}