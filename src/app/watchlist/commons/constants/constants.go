package constants

//Authentications API URL Keys
const (
	ServiceName       = "authentication"
	PortDefaultValude = 5030
)

// Database table name & field names for users
const (
	UsersTableName = "users"
	Fieldemail     = "email"
)

const (
	UsernameField = "username = ?"
)

// Success message for user login
const (
	UserCreationSuccessMsg = "User created successfully"
	UserLoggedInSuccessMsg = "User logged in successfully"
)

//success message for user otp generation
const (
	UserOtpGeneratedSuccess = "Otp generated successfully."
	UserOtpExpiryMsg        = "Otp will expire in 2 minutes."
)

//success message for user otp validation
const OtpValidatedSuccessMsg = "OTP validated successfully"

//success message for password change
const PasswordChangedSuccess = "Password changed successfully"

//Swagger Titile
const SwaggerTitle = "Stock Broker Application API"

const EmailorPasswordField = "email_or_password"

//Cookies
const (
	Name     = "refresh_token"
	Time     = 30 * 24 * 60 * 60
	Path     = "/"
	Domain   = ""
	Secure   = true
	HttpOnly = true
)

//purpose for JWT
const Purpose = "password_reset"

const ActionTypeSuccess = "success"
const ActionTypeFailure = "failure"



//queries
const WatchlistValidateQuery = `
	WITH input_watchlists AS (
		SELECT UNNEST(?::bigint[]) AS watchlist_id
	),

	scrip_exists AS (
		SELECT COUNT(*) > 0 AS exists_flag
		FROM scrip_masters
		WHERE id = ?
	),

	watchlist_user_check AS (
		SELECT id AS watchlist_id, user_id
		FROM watchlists
		WHERE id IN (SELECT watchlist_id FROM input_watchlists)
	),

	scrip_count_check AS (
		SELECT watchlist_id,
			COUNT(*) AS total_scrips
		FROM watchlist_scrips
		GROUP BY watchlist_id
	),

	scrip_present_check AS (
		SELECT watchlist_id
		FROM watchlist_scrips
		WHERE LOWER(scrip_id) = LOWER(?)
		AND watchlist_id IN (SELECT watchlist_id FROM input_watchlists)
	)

	SELECT
		i.watchlist_id AS "WatchlistId",

		-- ADD
		CASE
			WHEN NOT (SELECT exists_flag FROM scrip_exists) THEN false
			WHEN wu.watchlist_id IS NULL THEN false
			WHEN wu.user_id != ? THEN false
			WHEN COALESCE(sc.total_scrips,0) >= 10 THEN false
			WHEN sp.watchlist_id IS NOT NULL THEN false
			ELSE true
		END AS "CanInsert",

		-- DELETE
		CASE
			WHEN wu.watchlist_id IS NULL THEN false
			WHEN wu.user_id != ? THEN false
			WHEN sp.watchlist_id IS NULL THEN false
			ELSE true
		END AS "CanDelete",

		CASE
			WHEN wu.watchlist_id IS NULL
				THEN ?
			WHEN wu.user_id != ?
				THEN ?
			WHEN sp.watchlist_id IS NOT NULL
				AND COALESCE(sc.total_scrips,0) < 10
				THEN ?
			WHEN sp.watchlist_id IS NULL
				THEN ?
			WHEN COALESCE(sc.total_scrips,0) >= 10
				THEN ?
			ELSE 'OK'
		END AS "Reason"

	FROM input_watchlists i
	LEFT JOIN watchlist_user_check wu
		ON wu.watchlist_id = i.watchlist_id
	LEFT JOIN scrip_count_check sc
		ON sc.watchlist_id = i.watchlist_id
	LEFT JOIN scrip_present_check sp
		ON sp.watchlist_id = i.watchlist_id;
	`


const WatchlistAddQuery =
`
		WITH inserted AS (
			INSERT INTO watchlist_scrips (watchlist_id, scrip_id)
			SELECT id, ?
			FROM watchlists
			WHERE id IN (?)
			RETURNING watchlist_id
		)

		UPDATE watchlists w
		SET
			scrip_count = w.scrip_count + 1,
			last_updated_at = NOW()
		FROM inserted i
		WHERE w.id = i.watchlist_id
		RETURNING w.id AS "WatchlistId", w.watchlist_name AS "WatchlistName";
	`


const WatchlistDeleteQuery =
`
	WITH delete_from_watchlist_scrips AS(
			delete from watchlist_scrips
			where watchlist_id in (?) and scrip_id = ?
			returning watchlist_id
	),
	update_watchlists_table as(
			update watchlists
			set
			scrip_count = Greatest(scrip_count-1,0),
			last_updated_at= now()
			where id in (select watchlist_id from delete_from_watchlist_scrips)
	)
	SELECT 
		w.id AS watchlist_id,
		w.watchlist_name
	FROM watchlists w
	JOIN delete_from_watchlist_scrips d ON w.id = d.watchlist_id;
	`