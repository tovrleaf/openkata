import React from 'react';

interface ModalProps {
  title: string;
  children: React.ReactNode;
  onClose: () => void;
}

export const Modal: React.FC<ModalProps> = ({ title, children, onClose }) => (
  <div className="modal-overlay">
    <div className="modal">
      <h2>{title}</h2>
      <div className="modal-body">{children}</div>
      <button onClick={onClose}>Close</button>
    </div>
  </div>
);
