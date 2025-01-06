import { useState, useEffect } from "react";
import { useParams, useNavigate } from 'react-router-dom';
import axios from "axios";

const Transactions = () => {
    const redirect = useNavigate()
    const [sort, setSort] = useState()
    const {item_name, budget_name} = useParams()
    const [transactions, setTransactions] = useState([
            {
                transaction_id: "",
                transaction_name: "",
                transaction_type: "",
                amount: 0.00,
                date: ""
            }
        ])

    const [date, setDate] = useState(new Date())
    const token = localStorage.getItem("token")

    useEffect(() => {
        if (token === null) {
            redirect('/login')
        }
        const getTransactions = async () => {
            await axios.get(`http://localhost:8080/main/transactions/${budget_name}/${item_name}/${date.getFullYear()}/${date.getMonth() + 1}`, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            }).then(response => {
                setTransactions(response.data.transactions)
            }).catch(error => {
                alert(error.response?.data || error.message)
            })
        }

        getTransactions();
    }, [redirect, token, budget_name, item_name, date])
    
    const handleSort = (e) => {
        setSort(e.target.value)
        switch (sort) {

        }
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
                    <h1 className="px-5">Month: {date.getMonth() + 1}</h1>
                    <h1 className="px-5">Year: {date.getFullYear()}</h1>
                </div>
                <table>
                    <tr>
                        <th className="px-5">Name</th>
                        <th className="px-5">Type</th>
                        <th className="px-5">Amount</th>
                        <th className="px-5">Date</th>
                    </tr>
                    {/* Array.isArray(b.budget_items) && b.budget_items[0] !== null ? */}
                    {transactions.map((transaction, index) => (
                    <tr key={index}>
                        <td className="px-5">{transaction.transaction_name}</td>
                        <td className="px-5">{transaction.transaction_type}</td>
                        <td className="px-5">{transaction.amount}</td>
                        <td className="px-5">{`${transaction.date.substring(0, 10)}`}</td>
                    </tr>
                    ))}
                </table>             
            </div>

        </div>
    )
}

export default Transactions;