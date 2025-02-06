<template>
    <div class="row g-0">

        <div class="col-md-1 col-2">

            <slot name="conversationImage"></slot>

        </div>

        <div class="col-md-11 col-10">
            <div class="card-body">

                <slot name="conversationName"></slot>

                <div class="d-flex justify-content-between">
                    <slot name="conversationMessage"></slot>
                </div>

            </div>
        </div>

    </div>
</template>