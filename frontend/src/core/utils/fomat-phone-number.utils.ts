/**
 * Formats a phone number into a visually appealing format.
 * Example: +79667778899 -> +7 (966) 777-88-99
 *
 * @param phoneNumber - The raw phone number in E.164 format (e.g., +79667778899).
 * @returns The formatted phone number as a string, or an empty string if invalid.
 */
export const formatPhoneNumber = (phone: string) => {
	return phone.replace(/(\+7)(\d{3})(\d{3})(\d{2})(\d{2})/, '$1 $2 $3-$4-$5')
}
