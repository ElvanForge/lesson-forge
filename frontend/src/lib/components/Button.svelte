<script lang="ts">
    interface ButtonProps {
        text: string;
        variant?: "primary" | "secondary" | "accent";
        onclick?: () => void;
        isLoading?: boolean;
        disabled?: boolean;
    }
    let {
        text,
        variant = "primary",
        onclick,
        isLoading = false,
        disabled = false,
    }: ButtonProps = $props();
</script>

<button
    {onclick}
    disabled={isLoading || disabled}
    class="btn-premium disabled:opacity-50 disabled:cursor-not-allowed
    {variant === 'primary' ? 'bg-primary text-white hover:bg-primary/90' : ''}
    {variant === 'secondary'
        ? 'bg-secondary text-white hover:bg-secondary/90'
        : ''}
    {variant === 'accent' ? 'bg-accent text-primary hover:bg-accent/90' : ''}"
>
    {#if isLoading}
        <svg
            class="animate-spin h-5 w-5 text-current"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
        >
            <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
            ></circle>
            <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
        </svg>
    {/if}
    {text}
</button>
