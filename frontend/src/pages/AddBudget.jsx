import { useEffect, useState } from "react";
import axios from "axios";
import { FaMinus } from "react-icons/fa";
import { useNavigate } from "react-router-dom";


const AddBudget = () => {
    
    const redirect = useNavigate();
    const [budgetName, setBudgetName] = useState("")
    const token = localStorage.getItem("token")
    const [items, setItems] = useState([{
            item_name: "",
            budget_cost: 0.00,
            description: "",
            priority: 0.00
    },])

    useEffect(() => {
        if (token === null) {
            redirect('/login')
        }
    }, [token, redirect])

    const handleAddBudget = async () => {
        // prepare budget details
        const details = {budget_name: budgetName, items: items}
        // send request to add budget to the database
        await axios.post("http://localhost:8080/main/budgets/add_budget", details, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        }).then(response => {
            // show success message
            alert(response.data.Message)
            // redirect to budgets
            redirect('/budgets')
        }).catch(error => { // show error message
            alert(error.response?.data || error.message);
        })
    }

    const handleBudgetChange = (e) => {
        const {value} = e.target;
        setBudgetName(value)
    }
    
    const handleItemChange = (e, index) => {
        const {name, value} = e.target;
        setItems((previousItems) =>
            previousItems.map((item, i) => 
                i === index ? 
                {...item, [name]: name === "priority" || name === "budget_cost" ? 
                    parseFloat(value) : value 
                } : item
        ))
    }

    const addItem = () => {
        setItems([...items, {
            item_name: "",
            budget_cost: 0.00,
            description: "",
            priority: 0.00
        }]);
      };

    const removeItem = (index) => {
        setItems((previousItems) => previousItems.filter((_, i) => i !== index));
    }

    return (
        <div className="p-20 grid">
            <div className="flex justify-evenly">
                <h1>Add Budget</h1>
            </div>
            <div className="grid">
                <input 
                type="text"
                name="budget_name"
                placeholder="Budget Name..."
                onChange={handleBudgetChange}
                className="py-2"
                required
                />
            </div>
            <div className="flex justify-evenly">
                <h1>Add Items</h1>
            </div>
            {items.map((item, index) => (
                <div key={index} className="grid">
                    <button onClick={() => removeItem(index)}><FaMinus /></button>
                    <input 
                    type="text"
                    name="item_name"
                    placeholder="Item Name..."
                    value={item.item_name}
                    onChange={(e) => handleItemChange(e, index)}
                    className="py-2"
                    required
                    />
                    <textarea 
                    type="text"
                    name="description"
                    placeholder="Description..."
                    value={item.description}

                    onChange={(e) => handleItemChange(e, index)}
                    className="py-2"
                    required
                    />
                    <div className="flex items-center">
                        <h1 className="px-2">£</h1>
                        <input 
                        type="number"
                        name="budget_cost"
                        placeholder="Budget Amount..."
                        value={item.budget_cost}
                        onChange={(e) => handleItemChange(e, index)}
                        className="p-2"
                        required
                        />                        
                    </div>

                    <h1>Priority</h1>
                    <input 
                    type="range"
                    min="1"
                    max="10"
                    step="1"
                    name="priority"
                    value={item.priority}
                    onChange={(e) => handleItemChange(e, index)}
                    className="py-2"
                    required
                    />
                </div>
            ))}
            <button onClick={addItem}>Add another item</button>
            <button onClick={handleAddBudget}>Create Budget</button>
            <button onClick={(e) => {e.preventDefault();redirect('/budgets')}}>Back to my budget</button>
        </div>
    )
}

export default AddBudget;