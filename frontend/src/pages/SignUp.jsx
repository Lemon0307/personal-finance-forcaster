import React, { useState } from "react";
import { useNavigate } from 'react-router-dom'; 
import axios from "axios";

const SignUp = () => {
  // Step management
  const [step, setStep] = useState(1);

  // User information state
  const [details, setDetails] = useState({
    username: "",
    email: "",
    password: "",
    confirm_password: "",
    forename: "",
    surname: "",
    dob: "",
    address: "",
    current_balance: "",
  });

  // Redirect function
  let redirect = useNavigate()

  // Security questions state
  const [securityQuestions, setSecurityQuestions] = useState([
    { question: "", answer: "" },
  ]);

  // Available security questions
  const questions = [
    "What is your favorite color?",
    "What is your mother's maiden name?",
    "What was your first pet's name?",
    "What city were you born in?",
  ];

  const handleDetailsChange = (e) => {
    const { name, value } = e.target;
    setDetails({ ...details, [name]: value });
  };

  const handleSecurityQuestionChange = (index, field, value) => {
    const updatedQuestions = [...securityQuestions];
    updatedQuestions[index][field] = value;
    setSecurityQuestions(updatedQuestions);
  };

  const addSecurityQuestion = () => {
    setSecurityQuestions([...securityQuestions, { question: "", answer: "" }]);
  };

  console.log(details)


  const handleSubmit = async (e) => {
    try {
      const userData = {
        user: {...details, current_balance: parseFloat(details.current_balance)},
        security_questions: securityQuestions,
      };
      console.log(userData)
    const response = await axios.post("http://localhost:8080/sign_up", userData)
    alert(response.data.Message)
    e.preventDefault()
    // redirects to login page after success message
    redirect('/login')
    } catch (error) {
        alert(error.response.data)
    }
  };

  return (
    <div className="p-40">
      {step === 1 && (
        <div className="grid w-11/12 justify-center">
          <h2>User Information</h2>
          <input
            type="text"
            name="username"
            placeholder="Username"
            value={details.username}
            onChange={handleDetailsChange}
          />
          <input
            type="email"
            name="email"
            placeholder="Email"
            value={details.email}
            onChange={handleDetailsChange}
          />
          <input
            type="password"
            name="password"
            placeholder="Password"
            value={details.password}
            onChange={handleDetailsChange}
          />
          <input
            type="password"
            name="confirm_password"
            placeholder="Confirm Password"
            value={details.confirm_password}
            onChange={handleDetailsChange}
          />
          <input
            type="text"
            name="forename"
            placeholder="Forename"
            value={details.forename}
            onChange={handleDetailsChange}
          />
          <input
            type="text"
            name="surname"
            placeholder="Surname"
            value={details.surname}
            onChange={handleDetailsChange}
          />
          <input
            type="date"
            name="dob"
            placeholder="Date of Birth"
            value={details.dob}
            onChange={handleDetailsChange}
          />
          <textarea
            name="address"
            placeholder="Address"
            value={details.address}
            onChange={handleDetailsChange}
          />
          <input
            type="number"
            name="current_balance"
            placeholder="Current Balance"
            value={details.current_balance}
            onChange={handleDetailsChange}
          />
          <button onClick={() => setStep(2)}>Next</button>
        </div>
      )}

      {step === 2 && (
        <div>
          <h2>Security Questions</h2>
          {securityQuestions.map((sq, index) => (
            <div key={index}>
              <select
                value={sq.question}
                onChange={(e) =>
                  handleSecurityQuestionChange(index, "question", e.target.value)
                }
              >
                <option value="" disabled>
                  Select a question
                </option>
                {questions.map((question, i) => (
                  <option key={i} value={question}>
                    {question}
                  </option>
                ))}
              </select>
              <input
                type="text"
                placeholder="Answer"
                value={sq.answer}
                onChange={(e) =>
                  handleSecurityQuestionChange(index, "answer", e.target.value)
                }
              />
            </div>
          ))}
          <button onClick={addSecurityQuestion}>Add Another Question</button>
          {/* <button onClick={() => setStep(1)}>Back</button> */}
          <button onClick={handleSubmit}>Submit</button>
        </div>
      )}
    </div>
  );
};

export default SignUp;
