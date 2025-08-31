interface Props {
  children: React.ReactNode;
  disabled?: boolean;
  onClick?: () => void;
  style?: React.CSSProperties;
  variant?: "pink" | "blue" | "red" | "green";
  size?: "sm" | "md" | "lg";
}
export const Button = ({
  disabled,
  onClick,
  style,
  variant = "pink",
  size = "md",
  children,
}: Props) => {
  const sizeStyles = {
    sm: { height: "24px", padding: "0 8px", fontSize: "12px" },
    md: { height: "32px", padding: "0 12px", fontSize: "14px" },
    lg: { height: "48px", padding: "0 16px", fontSize: "16px" },
  };

  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled}
      className="button-component"
      style={{
        borderColor: disabled ? "var(--gray-400)" : `var(--${variant}-600)`,
        backgroundColor: disabled ? "var(--gray-300)" : `var(--${variant}-500)`,
        color: disabled ? "var(--gray-500)" : "white",
        height: `${sizeStyles[size].height}`,
        padding: `${sizeStyles[size].padding}`,
        fontSize: `${sizeStyles[size].fontSize}`,
        borderRadius: "6px",
        borderWidth: "1px",
        borderStyle: "solid",
        cursor: disabled ? "not-allowed" : "pointer",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        gap: "8px",
        transform: "scale(1)",
        transition:
          "transform 50ms, background-color 200ms, border-color 200ms, box-shadow 200ms",
        ...style,
      }}
      onMouseEnter={(e) => {
        if (!disabled) {
          e.currentTarget.style.boxShadow =
            `0 0 8px var(--${variant}-100), ` +
            `0 0 16px var(--${variant}-500)`;
        }
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.boxShadow = "none";
      }}
    >
      {children}
    </button>
  );
};
