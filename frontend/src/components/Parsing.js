export const parseCSVToJSON = (file) => {
    const lines = file.split("\n")
    const transactions = []
    console.log(lines)

    for (let i = 2; i < lines.length; i++) {
        const row = lines[i].trim().replace(/^,/, "")
        const col = row.split(",").map(i => i.trim())
        const transaction = {
            date: parseDate(col[0]),
            name: col[1],
            amount: parseFloat(col[2]),
            type: col[3].toLowerCase()
        }
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