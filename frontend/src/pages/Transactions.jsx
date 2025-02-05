import { useState, useEffect } from "react"
import { useParams, useNavigate } from 'react-router-dom'
import axios from "axios"
import { FaPlus, FaMinus} from 'react-icons/fa'
import { quickSort } from "../components"


const Transactions = () => {
    const redirect = useNavigate()
    const [sort, setSort] = useState()
    const {item_name, budget_name} = useParams()
    const [transactions, setTransactions] = useState([
        {
            transaction_id: "",
            name: "",
            type: "",
            amount: 0.00,
            date: ""
        }
    ])

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
            console.log(date)
            await axios.get(`http://localhost:8080/main/transactions/${budget_name}/${item_name}/${date.year}/${date.month}`, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            }).then(response => {
                setTransactions(response.data.transactions)
                console.log(response.data.transactions)
            }).catch(error => {
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

    const handleSubmit = async () => {
        let ok = true
        for (const key in newTransaction) {
            if (typeof newTransaction[key] === "string" && newTransaction[key].trim().length === 0) {
                ok = false;
            }
        }

        if (!ok) {
            alert("Please fill in all the required details.");
            return;
        }
        const reqData = {
            item: {
                budget_name: budget_name,
                item_name: item_name
            }, transactions: [{
                    ...newTransaction,
                    amount: parseFloat(newTransaction.amount)                           
            }]
        }
        await axios.post("http://localhost:8080/main/transactions/add_transaction", reqData, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        }).then(response => {
            console.log(response.data.Message)
            window.location.reload()
        }).catch(error => {
            alert(error.response?.data || error.message)
        })
    }

    const handleRemoveTransaction = async (transaction_id) => {
        console.log(transaction_id)
        await axios.delete(`http://localhost:8080/main/transactions/${date.year}/${date.month}/${budget_name}/${item_name}/remove_transaction/${transaction_id}`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
        .then(response => {
            console.log(response.data.Message)
            window.location.reload() 
        }).catch(error => {
            alert(error.response?.data || error.message)
        })
    }
    
    const handleSort = (e) => {
        setTransactions((previousTransactions) => (
            previousTransactions = quickSort(previousTransactions, e.target.value)
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

    return (
        <div className="p-20 flex">
            <div className="px-10">
                <h1>Sort by:</h1>
                <select value={sort} onChange={handleSort}>
                    <option value="name">Name</option>
                    <option value="type">Type</option>
                    <option value="amount">Amount</option>
                    <option value="date">Date</option>
                </select>
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
                        <button onClick={() => handleSubmit()}>
                            <FaPlus color="grey"/>    
                        </button>
                    </td>
                    </tr>
                </table>             
            </div>

        </div>
    )
}

export default Transactions;