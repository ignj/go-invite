export interface Celebration {
  _id: string;
  date: string;
  title: string;
  description: string;
  attendees: Attendee[];
}

export interface Attendee {
  _id: string;
  fullName: string;
  email: string;
  phoneNumber: string;
  attendStatus: number;
}
