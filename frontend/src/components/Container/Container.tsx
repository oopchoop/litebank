import React from "react";
import styles from "./Container.module.css"

interface ContainerProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
}

export const Container: React.FC<ContainerProps> = ({ children, className, ...props }) => {
  const combinedClasses = `${styles.container} ${className || ''}`;
  
  return (
    <div className={combinedClasses} {...props}>
      {children}
    </div>
  );
};

