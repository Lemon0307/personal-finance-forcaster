import { useNavigate } from "react-router-dom";

const Navbar = () => {
    const redirect = useNavigate()

    const logoutHandler = (e) => {
        localStorage.removeItem('token')
        localStorage.removeItem('email')
        localStorage.removeItem('username')
        redirect('/login')
    }

    return (
        <div className="flex justify-around">
            <h1>Personal Finance Forecaster</h1>
            <div>
                <button className="px-10" onClick={(e) => {e.preventDefault(); redirect('/')}}>Home</button>
                <button className="px-10" onClick={(e) => {logoutHandler()}}>Logout</button>                
            </div>

        </div>
    )
}

export default Navbar;