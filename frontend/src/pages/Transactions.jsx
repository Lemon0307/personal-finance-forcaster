import { useState, useEffect, useRef } from "react"
import { useParams, useNavigate } from 'react-router-dom'
import axios from "axios"
import { FaPlus, FaMinus} from 'react-icons/fa'
import { quickSort } from "../components"
import { parseCSVToJSON } from "../components/Parsing.js"


const Transactions = () => {
    const redirect = useNavigate()
    // sort options
    const [sort, setSort] = useState("name")
    const [order, setOrder] = useState("asc")
    // stores the item name and budget name of transactions
    const {item_name, budget_name} = useParams()
    // stores the list of transactions by the user
    const [transactions, setTransactions] = useState([
        {
            transaction_id: "",
            date: "",
            type: "",
            name: "",
            amount: 0.00,

        }
    ])
    // stores the csv file to be imported
    const [csvFile, setCSVfile] = useState()

    // stores data of new transaction to be added
    const [newTransaction, setNewTransaction] = useState({
        name: "",
        type: "",
        amount: 0.00,
        date: ""
    })

    const [date, setDate] = useState({
        month: new Date().getMonth() + 1,
        year: new Date().getFullYear()
    })

    const [dateString, setDateString] = useState(() => {
        const date_now = new Date()
        return `${date_now.getFullYear()}-${String(date_now.getMonth() + 1).padStart(2, "0")}`
    })

    const token = localStorage.getItem("token")

    useEffect(() => {
        if (token === null) {
            redirect('/login')
        }
        const getTransactions = async () => {
            // send request to get transactions from the user of the current month year
            await axios.get(`http://localhost:8080/main/transactions/${budget_name}/${item_name}/${date.year}/${date.month}`, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            }).then(response => {
                // store response data in the transactions state
                setTransactions(response.data.transactions)
            }).catch(error => {
                // show error message
                alert(error.response?.data || error.message)
            })
        }

        getTransactions();
    }, [redirect, token, budget_name, item_name, date])

    const handleTransactionChange = (e) => {
        const {name, value} = e.target
        setNewTransaction(previousTransaction => ({
            ...previousTransaction,
            [name]: value
        }))
    }

    const handleAddTransaction = async () => {
        let ok = true
        for (const key in newTransaction) { // checks if each required detail is entered
            if (typeof newTransaction[key] === "string" && newTransaction[key].trim().length === 0) {
                ok = false;
            }
        }

        if (!ok) { // shows error message if not all details are entered
            alert("Please fill in all the required details.");
            return;
        }
        // gather all transaction data
        const requestData = {
            item: {
                budget_name: budget_name,
                item_name: item_name
            }, transactions: [{
                    ...newTransaction,
                    amount: parseFloat(newTransaction.amount)                           
            }]
        }
        // sends request to add a transaction to the database
        await axios.post("http://localhost:8080/main/transactions/add_transaction", requestData, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        }).then( // appends the new transaction to the transactions state
            setTransactions(prev => Array.isArray(prev) ? [...prev, newTransaction] : [newTransaction]) 
        ).catch(error => {
            // show error message
            alert(error.response?.data || error.message)
        })
    }

    const handleRemoveTransaction = async (transaction_id) => {
        // send request to delete transaction from the database
        await axios.delete(`http://localhost:8080/main/transactions/${date.year}/${date.month}/${budget_name}/${item_name}/remove_transaction/${transaction_id}`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
        .then( // removes transaction from transactions state
            setTransactions(prev => prev.filter(transaction => transaction.transaction_id !== transaction_id))
        ).catch(error => {
            alert(error.response?.data || error.message)
        })
    }
    
    const handleSort = () => {
        setTransactions((previousTransactions) => (
            // run quicksort on transactions
            previousTransactions = quickSort(previousTransactions, sort, order)
        ));
    }

    const handleMonthYearChange = (e) => {
        const date_string = e.target.value.split('-')
        setDateString(e.target.value)
        setDate(previousDate => ({
            ...previousDate,
            month: parseInt(date_string[1]),
            year: parseInt(date_string[0])
        }))
    }

    const handleExportCSV = () => {
        // does nothing if there are no transactions
        if (!transactions) {
            return;
        }
        // sets the headers of the csv as the keys of the transactions
        const headers = Object.keys(transactions[0]).slice(1)
        const rows = transactions.map(transaction => {
            return headers.map(header => {
                let value = transaction[header];
    
                // Check if the current field is a date and format it
                if (header.toLowerCase().includes("date") && value) {
                    const date = new Date(value);
                    if (!isNaN(date.getTime())) {
                        value = date.toLocaleDateString("en-GB"); // Converts to dd/mm/yyyy format
                    }
                }
    
                return value;
            }).join(",");
        })
        // combines headers and rows to make the file
        const file = [headers.join(","), ...rows].join("\n")
        const blob = new Blob([file], {type: "text/csv"})
        // download the file for the user
        const link = document.createElement("a")
        link.href = URL.createObjectURL(blob)
        // file name is MM-YYYY-item_name-budget_name
        link.download = `${date.month}-${date.year}-${item_name}-${budget_name}`
        link.click()
    }

    const handleImportCSV = () => {
        const reader = new FileReader()
        
        // read the csv file imported
        reader.readAsText(csvFile);
        reader.onload = async () => {
            const csv = reader.result
            console.log(csv)
            const transactions = parseCSVToJSON(csv)

            // gather import data
            const importData = {
                item: {
                    budget_name: budget_name,
                    item_name: item_name,
                },
                transactions: transactions,
                monthly_costs: {
                    month: new Date().getMonth() + 1,
                    year: new Date().getFullYear()
                }
            }
            // send request to add the new imported transactions to the database
            await axios.post(`http://localhost:8080/main/transactions/add_transaction`, 
                importData,
                {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            }).then(response => { // show success message
                alert(response.data.Message)
                // refresh screen
                window.location.reload()
            }).catch(error => { // show error message
                alert(error.response?.data)     
            })
        }
    }

    return (
        <div className="p-20 flex">
            <div className="px-10">
                <h1>Sort by:</h1>
                <select value={sort.current} onChange={(e) => setSort(e.target.value)}>
                    <option value="name">Name</option>
                    <option value="type">Type</option>
                    <option value="amount">Amount</option>
                    <option value="date">Date</option>
                </select>
                <select value={order.current} onChange={(e) => setOrder(e.target.value)}>
                    <option value="asc">Ascending</option>
                    <option value="desc">Descending</option>
                </select>
                <button onClick={handleSort}>Sort</button>
            </div>
            <div className="grid place-content-center">
                <div className="flex justify-evenly">
                    <h1 className="px-5">My Transactions</h1>
                    <h1 className="px-5">Item: {item_name}</h1>
                    <h1 className="px-5">Budget: {budget_name}</h1>
                    <input type="month" value={dateString} onChange={(e) => {
                        handleMonthYearChange(e);
                    }}/>
                </div>
                <table>
                    <tr>
                        <th className="px-5">Name</th>
                        <th className="px-5">Type</th>
                        <th className="px-5">Amount</th>
                        <th className="px-5">Date</th>
                    </tr>

                    {Array.isArray(transactions) && transactions[0] !== null ? 
                    transactions.map((transaction, index) => (
                    <tr key={index}>
                        <td className="px-5">{transaction.name}</td>
                        <td className="px-5">{transaction.type}</td>
                        <td className="px-5">£{transaction.amount}</td>
                        <td className="px-5">{`${transaction.date.substring(0, 10)}`}</td>
                        <td>
                        <button onClick={
                            () => handleRemoveTransaction(transaction.transaction_id)}>
                            <FaMinus color="grey"/>    
                        </button>
                    </td>
                    </tr>
                    )): <div>There are no transactions for this budget item</div>}
                    <tr>
                    <td className="px-5" >
                        <input type="text" 
                        name="name" 
                        placeholder="Name..." 
                        required 
                        value={newTransaction.name}
                        onChange={(e) => handleTransactionChange(e)}/>
                    </td>
                    <td className="px-5" >
                        <select name="type" 
                        onChange={(e) => {handleTransactionChange(e)}
                        }
                        className="p-2 border-r- bg-zinc-100"
                        >
                            <option value="" selected>Transaction type...</option>
                            <option value="inflow">Inflow (+)</option>
                            <option value="outflow">Outflow (-)</option>
                        </select>
                    </td>
                    <td className="px-5" >
                        <input type="number" 
                        name="amount" 
                        placeholder="Amount..." 
                        required 
                        value={newTransaction.amount}
                        onChange={(e) => handleTransactionChange(e)}/>
                    </td>
                    <td className="px-5" >
                        <input type="date" 
                        name="date" 
                        placeholder="Date..." 
                        required 
                        value={newTransaction.date}
                        onChange={(e) => handleTransactionChange(e)}/>
                    </td>
                    <td>
                        <button onClick={() => handleAddTransaction()}>
                            <FaPlus color="grey"/>    
                        </button>
                    </td>
                    </tr>
                </table>             
            </div>
            <div>
                <button className="p-5" onClick={() => handleExportCSV()}>Export to CSV</button>
                <input type="file" accept=".csv" onChange={(e) => setCSVfile(e.target.files[0])}/>
                <button className="p-5" onClick={() => handleImportCSV()}>Import from CSV</button>                
            </div>

        </div>
    )
}

export default Transactions;