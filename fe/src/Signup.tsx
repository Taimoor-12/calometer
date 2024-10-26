import { useState } from "react";
import { Link } from "react-router-dom";
import "./Signup.css";

interface ValidationErrors {
  fullName?: string;
  username?: string;
  password?: string;
  confirmPassword?: string;
}

function Signup() {
  const [formData, setFormData] = useState({
    fullName: "",
    username: "",
    password: "",
    confirmPassword: "",
  });

  const [errors, setErrors] = useState<ValidationErrors>({})

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData({ ...formData, [name]: value })

    const error = validateField(name, value)
    setErrors({ ...errors, [name]: error })
  }

  const validateField = (name: string, value: string): string | undefined => {
    switch(name) {
      case "fullName":
        if (!value.trim()) {
          return "Full name is required"
        } else if (!/^[A-Za-z\s]+$/.test(value)) {
          return "Full name must contain only letters"
        }
        break
      case "username":
        if (!value.trim()) {
          return "Username is required"
        }
        break
      case "password":
        if (!value.trim()) {
          return "Password is required"
        } else if (value.length < 6) {
          return "Password must be at least 6 characters"
        }
        break
      case "confirmPassword":
        if (!value.trim()) {
          return "Confirm password is required"
        } else if (value !== formData.password) {
          return "Passwords do no match"
        }
        break
      default:
        break
    }

    return undefined
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    const validationErrors: ValidationErrors = {}

    if (!formData.fullName.trim()) {
      validationErrors.fullName = "Full name is required"
    } else if (!/^[A-Za-z\s]+$/.test(formData.fullName)) {
      validationErrors.fullName = "Full name must contain only letters"
    }

    if (!formData.username.trim()) {
      validationErrors.username = "Username is required"
    }

    if (!formData.password.trim()) {
      validationErrors.password = "Password is required"
    } else if (formData.password.length < 6) {
      validationErrors.password = "Password must be at least 6 characters"
    }

    if (!formData.confirmPassword.trim()) {
      validationErrors.confirmPassword = "Confirm password is required"
    } else if (formData.password !== formData.confirmPassword) {
      validationErrors.confirmPassword = "Passwords do not match"
    }

    setErrors(validationErrors)


    if(Object.keys(validationErrors).length === 0) {
      console.log("Form submitted", formData)
    }
  };

  return (
    <div className="main">
      <h1>SIGNUP</h1>

      <div className="parentDiv">
        <div className="formDiv">
          <form onSubmit={handleSubmit}>
            <div className="inputDiv">
              <p>Full Name</p>
              <input 
                type="text"
                name="fullName"
                value={formData.fullName} 
                placeholder="Enter your fullname"
                onChange={handleChange} 
              />
              { errors.fullName && <p className="error">{errors.fullName}</p> }
            </div>

            <div className="inputDiv">
              <p>Username</p>
              <input 
                type="text"
                name="username"
                value={formData.username} 
                placeholder="Enter your username"
                onChange={handleChange} 
              />
              { errors.username && <p className="error">{errors.username}</p> }
            </div>

            <div className="inputDiv">
              <p>Password</p>
              <input 
                type="password"
                name="password"
                value={formData.password} 
                placeholder="Enter your password"
                onChange={handleChange} 
              />
              { errors.password && <p className="error">{errors.password}</p> }
            </div>

            <div className="inputDiv">
              <p>Confirm Password</p>
              <input 
                type="password"
                name="confirmPassword"
                value={formData.confirmPassword} 
                placeholder="Enter the password again"
                onChange={handleChange} 
              />
              { errors.confirmPassword && <p className="error">{errors.confirmPassword}</p> }
            </div>
            
            <div className="registerBtn">
              <button type="submit">Register</button>
            </div>
          </form>
        
        <p className="loginPrompt">
          Already have an account? <span><Link to="/login">Login here</Link></span>
        </p>
        </div>
        <div className="imgDiv">
          <img src="/assets/healthy_habit.svg" alt="healthy habit" />
        </div>
      </div>
    </div>
  );
}

export default Signup;
