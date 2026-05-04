export interface User {
    id: string;
    user_name: string;
    full_name: string;
    email: string;
    initials: string;
    is_guest: boolean;
    venmo_handle?: string;
    zelle_handle?: string;
    color?: string;
    created_at: string;
}

export interface Team {
    id: string;
    city: string;
    mascot: string;
    sport: 'football' | 'basketball' | 'baseball' | 'hockey' | 'soccer';
    full_name: string;
    team_logo_url?: string | null;
    primary_color: string;
    secondary_color: string;
}

export interface Square {
    id: string;
    board_id: string;
    row_index: number;
    col_index: number;
    user_id: string | null;
    user_color?: string;    // Added for live preview
    user_initials?: string; // Added for initials display
    first_name?: string; // Added for Player List
    last_name?: string;  // Added for Player List
    is_paid: boolean;
    payment_status: number; // 0 = Unpaid, 1 = Paid, 2 = Pending
    claimed_at: string | null;
}

export interface SquareUpdate {
    row: number; // Go: json:"row_index"
    col: number; // Go: json:"col_index"
    user_id: string;
    user_color: string;
    user_initials: string;
    first_name: string;
    last_name: string;
    is_paid: boolean;
    payment_status: number;
}

export interface Board {
    id: string;
    title: string;
    status: string;
    home_axis_numbers: number[];
    away_axis_numbers: number[];
    owner_id: string;
    price_per_square: number;
    creator_id: string;
}

export interface QuarterScore {
    quarter: number;
    home_score: number;
    away_score: number;
    is_locked: boolean;
}

export interface BoardScores {
    home: number;
    away: number;
    quarters: QuarterScore[]; // This is the missing piece!
}

export interface SquareOwnerDetails {
    first_name: string;
    last_name: string;
    color: string;
}

export interface BoardSummary {
    id: string;
    title: string;
    status: string;
    price_per_square: number;
    creator_id: string;
    home_team: Team; 
    away_team: Team;
    created_at?: string;
}

// Ensure your components use this specific shape for the data prop
export interface BoardPageData {
    board: Board;
    squares: Square[][];
    homeTeam: Team;
    awayTeam: Team;
    user: User | null;
    isCreator: boolean;
    token: string;
    payouts: Payouts[];
    scores: BoardScores;
    square_owners: Record<string, SquareOwnerDetails>;
    creator_venmo: string;
    creator_zelle: string;
}

export interface BoardPlayer {
    user_id: string;
    first_name: string;
    last_name: string;
    color: string;
    square_count: number;
    payment_status: number;
    is_paid: boolean; //TODO: remove this and just use payment_status
}

export interface SquareUpdate {
    row: number;
    col: number;
    user_id: string;
    user_color: string;
    user_initials: string;
}

export interface AxisUpdate {
    type: 'AXIS_UPDATE';
    home: number[];
    away: number[];
}

export type PayoutPeriod = 'q1' | 'q2' | 'q3' | 'q4' | 'final';

export interface Payouts {
    id: string;
    board_id: string;
    period_name: PayoutPeriod;
    amount: number;
    winner_user_id: string | null;
    winner_name: string | null;
    awarded_at: string | null; // ISO string from Go time.Time
}