import React from 'react';

export interface SelectOption {
  value: string;
  label: string;
}

interface SelectProps extends React.SelectHTMLAttributes<HTMLSelectElement> {
  label: string;
  options: string[] | SelectOption[];
  error?: string;
  placeholder?: string;
}

export const Select: React.FC<SelectProps> = ({ label, options, error, className = '', placeholder, ...props }) => {
  // Check if options are objects or strings
  const isObjectArray = options.length > 0 && typeof options[0] === 'object';

  return (
    <div className="mb-4">
      <label className="block text-gray-700 dark:text-gray-200 text-sm font-bold mb-2" htmlFor={props.id}>
        {label}
      </label>
      <select
        className={`shadow border rounded w-full py-2 px-3 text-gray-700 dark:text-gray-200 dark:bg-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent ${
          error ? 'border-red-500' : 'border-gray-300 dark:border-gray-600'
        } ${className}`}
        {...props}
      >
        <option value="">{placeholder || `Select ${label}`}</option>
        {isObjectArray
          ? (options as SelectOption[]).map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))
          : (options as string[]).map((option) => (
              <option key={option} value={option}>
                {option}
              </option>
            ))}
      </select>
      {error && <p className="text-red-500 text-xs italic mt-1">{error}</p>}
    </div>
  );
};
