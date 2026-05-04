export function getVenmoLink(
    handle: string, 
    amount: number, 
    boardTitle: string, 
    isMobile: boolean
): string {
    const cleanHandle = handle.replace('@', '');
    const note = encodeURIComponent(`Bellesquare: ${boardTitle}`);

    if (isMobile) {
        // Deep link for mobile apps
        return `venmo://paycharge?txn=pay&recipients=${cleanHandle}&amount=${amount}&notes=${note}`;
    }
    
    // Standard web URL for desktop
    // Note: Venmo's web interface doesn't support pre-filled notes/amounts via URL parameters
    return `https://venmo.com/u/${cleanHandle}`;
}