import { useEffect, useState } from "react";
import { FaPlus } from "react-icons/fa";
import axios from "axios"

const Budgets = () => {

    const [budgets, setBudgets] = useState([])
    const [budgetItem, setBudgetItem] = useState({
        item_name: "",
        budget_cost: 0.00,
        description: "",
        priority: 0.00
    })
    const token = localStorage.getItem('token')
    const [sort, setSort] = useState(null)

    useEffect(() => {
        const getBudgets = async () => {
            await axios.get("http://localhost:8080/main/budgets", {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            }).then(response => {
                setBudgets(response?.data)
            }).catch(error => {
                alert(error.response?.data || error.message)
            })
        }
        getBudgets()
    }, [token])
    console.log(budgets)

    const handleBudgetItemChange = (e) => {
        const {name, value} = e.target
        setBudgetItem({...budgetItem, [name]: value})
    }

    const handleSort = (e) => {
        setSort(e.target.value)
        switch (sort) {

        }
    }

    const handleSubmit = async (budget_name) => {
        console.log(budgetItem)
        budgetItem.priority = parseFloat(budgetItem.priority)
        budgetItem.budget_cost = parseFloat(budgetItem.budget_cost)
        console.log(typeof(budgetItem.budget_cost))
        await axios.post(`http://localhost:8080/main/budgets/add_budget_item/${budget_name}`, budgetItem, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        }).then(
            response => {
                alert(response.data.Message)
                window.location.reload();
            }
        ).catch(error => {
            alert(error.response.data)
        })
    }

    return (
        <div className="p-20 flex">
            <div className="px-10">
                <button>Add Budget</button>
                <h1>Sort by:</h1>
                <select value={sort} onChange={handleSort}>
                    <option value="item">Item</option>
                    <option value="amount">Amount</option>
                    <option value="priority">Priority</option>
                    <option value="description">Description</option>
                </select>
            </div>
            <div className="flex justify-evenly">
                <h1>My Budgets</h1>
            </div>
            <div className="grid place-content-center">
                {budgets.map((b, index) => (
                    <div key={index} className="py-5">
                            <div className="flex justify-around">
                                <h1>{b.budget.budget_name}</h1>
                                <div>
                                    <button className="px-2.5">Update Budget</button>
                                    <button className="px-2.5">Remove Budget</button>
                                </div>
                            </div>
                        <tbody>
                        <div key={index} className="flex justify-normal">
                            <table>
                            <tr>
                                <th className="px-5">Item</th>
                                <th className="px-5">Amount</th>
                                <th className="px-5">Priority</th>
                                <th className="px-5">Description</th>
                            </tr>
                            {b.budget_items.map((bi, index) => (
                                <tr key={index}>
                                  <td className="px-5">{bi.item_name}</td>
                                  <td className="px-5">{bi.budget_cost}</td>
                                  <td className="px-5">{bi.priority}</td>
                                  <td className="px-5">{bi.description}</td>
                                </tr>
                            ))}
                            <tr>
                                <td className="px-5" >
                                    <input type="text" name="item_name" placeholder="Item name..." onChange={handleBudgetItemChange}/>
                                </td>
                                <td className="px-5">
                                    <input type="number" name="budget_cost" placeholder="Amount..." onChange={handleBudgetItemChange}/>
                                </td>
                                <td className="px-5">
                                    <input type="number" name="priority" placeholder="Priority..." onChange={handleBudgetItemChange}/>
                                </td>
                                <td className="px-5">
                                    <input type="text" name="description" placeholder="Description..." onChange={handleBudgetItemChange}/>
                                </td>
                                <td>
                                    <button onClick={() => (handleSubmit(b.budget.budget_name))}>
                                        <FaPlus color="grey"/>    
                                    </button>
                                </td>
                            </tr>         
                            </table>
                        </div>
                        </tbody>
                    </div>
                ))}
            </div>
        </div>
    )
}

export default Budgets;