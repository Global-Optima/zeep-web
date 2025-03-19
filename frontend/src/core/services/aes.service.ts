import forge from 'node-forge'

export interface EncryptedData {
	iv: string // Base64-encoded IV (nonce)
	payload: string // Base64-encoded ciphertext + authentication tag
}

/**
 * Derives a 32-byte AES key from a given secret using SHA-256.
 */
function deriveAESKey(secret: string): string {
	return forge.util.bytesToHex(forge.sha256.create().update(secret).digest().bytes())
}

/**
 * Encrypts a JSON object using AES-GCM (256 bits).
 */
export function encryptPayload(data: unknown, secret: string): EncryptedData {
	// Convert object to JSON string
	const jsonString = JSON.stringify(data)

	// Generate a **12-byte IV** (Nonce)
	const iv = forge.random.getBytesSync(12)

	// Derive AES-256 key
	const keyHex = deriveAESKey(secret)
	const key = forge.util.hexToBytes(keyHex)

	// Encrypt using AES-GCM
	const cipher = forge.cipher.createCipher('AES-GCM', key)
	cipher.start({ iv })

	cipher.update(forge.util.createBuffer(jsonString, 'utf8'))
	cipher.finish() // Ensures authentication tag is generated

	// **Correctly extract authentication tag (16 bytes)**
	const authTag = cipher.mode.tag.bytes()
	const encryptedBytes = cipher.output.bytes() + authTag // Combine ciphertext + tag

	return {
		iv: forge.util.encode64(iv),
		payload: forge.util.encode64(encryptedBytes),
	}
}

/**
 * Decrypts an AES-GCM encrypted payload.
 */
export function decryptPayload(encryptedData: EncryptedData, secret: string): unknown {
	// Derive AES-256 key
	const keyHex = deriveAESKey(secret)
	const key = forge.util.hexToBytes(keyHex)

	// Decode Base64 IV and ciphertext
	const iv = forge.util.decode64(encryptedData.iv)
	const combinedCiphertext = forge.util.decode64(encryptedData.payload)

	// **Extract authentication tag (last 16 bytes)**
	const authTagLength = 16
	const ciphertext = combinedCiphertext.slice(0, -authTagLength)
	const authTag = combinedCiphertext.slice(-authTagLength)

	// Decrypt using AES-GCM
	const decipher = forge.cipher.createDecipher('AES-GCM', key)
	decipher.start({ iv, tagLength: 128 })
	decipher.update(forge.util.createBuffer(ciphertext))
	decipher.mode.tag = forge.util.createBuffer(authTag)

	const success = decipher.finish()

	if (!success) {
		throw new Error('Decryption failed: Invalid data or incorrect secret key.')
	}

	return JSON.parse(decipher.output.toString())
}
