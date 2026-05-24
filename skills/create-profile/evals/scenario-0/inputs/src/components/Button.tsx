import React from 'react';

interface ButtonProps {
  label: string;
  onClick: () => void;
  variant?: 'primary' | 'secondary';
}

export const Button: React.FC<ButtonProps> = ({ label, onClick, variant = 'primary' }) => (
  <button className={`btn btn-${variant}`} onClick={onClick}>
    {label}
  </button>
);
