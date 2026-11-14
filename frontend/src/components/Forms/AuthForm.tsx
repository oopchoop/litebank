import React from "react";
import styles from "./AuthForm.module.css";

interface AuthFormProps extends React.FormHTMLAttributes<HTMLFormElement> {
  header?: string;
  error?: string;
}

export const AuthForm: React.FC<React.PropsWithChildren<AuthFormProps>> = ({ header, children, error, ...props }) => {
  return (
    <form className={styles.container}
          {...props}>
      {header ? (
        <h3 className={styles.header}>{header}</h3>
      ) : (
        ''
      )}
      <span>{error}</span>
      {children}
    </form>
  )
}