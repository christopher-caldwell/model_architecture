use crate::router::auth::auth_middleware;
use crate::router::cors::get_cors;
use crate::router::dependencies::ServerDeps;
use axum::middleware::from_fn_with_state;
use axum::{
    routing::{get, post},
    Router,
};
use utoipa::OpenApi;
use utoipa_swagger_ui::SwaggerUi;

#[derive(OpenApi)]
#[openapi(
    nest(
        (path = crate::router::contact_inquiries::CONTACT_INQUIRIES_PATH, api = crate::router::contact_inquiries::ContactApi),
        (path = crate::router::health::HEALTH_CHECK_PATH, api = crate::router::health::HealthCheckApi)
    ),
    info(
        title = "Craft Code CRM",
        version = "1.0.0",
        description = "Internal management of Craft Code data",
        contact(name = "Support Team", email = "christopher@craftcode.solutions")
    )
)]
pub struct ApiDoc;

pub fn new_router(deps: ServerDeps) -> Router {
    let api = ApiDoc::openapi();

    let cors_layer = get_cors();

    let swagger_router = SwaggerUi::new("/swagger-ui").url("/api-docs/openapi.json", api);
    let public_router = Router::new()
        .route(
            crate::router::contact_inquiries::CONTACT_INQUIRIES_PATH,
            post(crate::router::contact_inquiries::create_contact_inquiry),
        )
        .route(
            crate::router::health::HEALTH_CHECK_PATH,
            get(crate::router::health::get_health_check),
        );

    let protected_router = Router::new()
        .route(
            crate::router::contact_inquiries::CONTACT_INQUIRIES_PATH,
            get(crate::router::contact_inquiries::get_contact_inquiries),
        )
        .route(
            crate::router::contact_inquiries::CONTACT_INQUIRY_BY_IDENT_PATH,
            get(crate::router::contact_inquiries::get_contact_inquiry_by_ident),
        )
        .layer(from_fn_with_state(deps.clone(), auth_middleware));

    Router::new()
        .merge(swagger_router)
        .merge(public_router)
        .merge(protected_router)
        .layer(cors_layer)
        .with_state(deps)
}
