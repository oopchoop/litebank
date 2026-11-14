import React from "react";
import styles from "./Button.module.css"

interface Button extends React.ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: 'primary' | 'danger'
}

export const Button: React.FC<Button> = ({variant = 'primary', className, children, ...props}) => {
    const classes = `${styles[variant]} ${className || ''}`
    
    return (
        <button {...props} className={classes}>
            {children}
        </button>
    )
}