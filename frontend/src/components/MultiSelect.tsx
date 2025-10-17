import React from 'react';

export interface MultiSelectOption {
  value: string;
  label: string;
}

interface MultiSelectProps {
  label: string;
  options: string[] | MultiSelectOption[];
  selected: string[];
  onChange: (selected: string[]) => void;
}

export const MultiSelect: React.FC<MultiSelectProps> = ({ label, options, selected, onChange }) => {
  // Check if options are objects or strings
  const isObjectArray = options.length > 0 && typeof options[0] === 'object';

  const handleToggle = (value: string) => {
    if (selected.includes(value)) {
      onChange(selected.filter((item) => item !== value));
    } else {
      onChange([...selected, value]);
    }
  };

  return (
    <div className="mb-4">
      <label className="block text-gray-700 dark:text-gray-200 text-sm font-bold mb-2">{label}</label>
      <div className="flex flex-wrap gap-2">
        {isObjectArray
          ? (options as MultiSelectOption[]).map((option) => (
              <button
                key={option.value}
                type="button"
                onClick={() => handleToggle(option.value)}
                className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                  selected.includes(option.value)
                    ? 'bg-primary-600 text-white'
                    : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-600'
                }`}
              >
                {option.label}
              </button>
            ))
          : (options as string[]).map((option) => (
              <button
                key={option}
                type="button"
                onClick={() => handleToggle(option)}
                className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                  selected.includes(option)
                    ? 'bg-primary-600 text-white'
                    : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-600'
                }`}
              >
                {option}
              </button>
            ))}
      </div>
    </div>
  );
};
